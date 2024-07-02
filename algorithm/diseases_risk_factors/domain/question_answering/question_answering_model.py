import logging
import os
import json
from typing import List
import collections
import aiohttp
import torch
import numpy as np
from sklearn.model_selection import train_test_split
from datasets import load_metric, Dataset, DatasetDict
from transformers import AutoTokenizer, AutoModelForQuestionAnswering, TrainingArguments, Trainer, default_data_collator
import pandas as pd
from tqdm.auto import tqdm
import evaluate
import torch.nn.functional as F
from random import shuffle
from tenacity import retry, stop_after_attempt, wait_fixed

from thesis_diseases_risk_factors.domain.question_answering._qa_dataset import _QADataset
from thesis_diseases_risk_factors.graph.client import Client, ListQuestionsAnswersByDiseaseQas


n_best = 20
"""
n_best is a parameter that specifies how many of the most probable start and end positions of the answer should be considered in the final stage of the model's predictions.
n_best is used to select the top n_best start and end positions with the highest probabilities (logits). These positions are then further processed, considering only answers 
that are in the context and with a length that is within a given range.
"""

class QuestionAnsweringModel:
    def __init__(self,
        graph_client: Client,
        http_client: aiohttp.ClientSession,
        trained_save_path: str) -> None:
        self._logger = logging.getLogger(__name__)
        self._graph_client = graph_client
        self._http_client = http_client
        self._trained_save_path = trained_save_path
        self._max_answer_length_full_file_name = os.path.join(self._trained_save_path, "max_answer_length.json")

    async def train_async(self, device: torch.device):
        diseases_response = await self._graph_client.list_diseases()

        # Get a list of all disease IDs
        all_disease_ids = [disease.id for disease in diseases_response.diseases]

        # Shuffle this list to randomize
        shuffle(all_disease_ids)

        # Perform an 80-20 split on disease IDs
        train_size = int(0.8 * len(all_disease_ids))
        train_disease_ids = all_disease_ids[:train_size]
        test_disease_ids = all_disease_ids[train_size:]

        # Create empty lists to hold training and testing data
        temp_train_data = []
        test_data = []

        for disease_id in all_disease_ids:
            # self._logger.info(f"Loading {disease.names[0]} QAs items")
            qas_response = await self._graph_client.list_questions_answers_by_disease(disease_id=disease_id)
            if qas_response.qas:
                for qa in qas_response.qas:
                    if len(qa.questions) == 0:
                        continue

                    entry = {"title": qa.article.id, 
                            "paragraphs": [{"context": qa.article.text, 
                                            "qas": []}]
                            }

                    for question in qa.questions:
                        if len(question.answers) == 0:
                            continue
                        qas_entry = {"id": question.id, 
                                    "question": question.text, 
                                    "answers": [],
                                    "is_impossible": False}
                        
                        for answer in question.answers:
                            qas_entry["answers"].append({"answer_start": answer.answer_start, 
                                                        "text": answer.text})
                            
                        entry["paragraphs"][0]["qas"].append(qas_entry)
                    
                    
                    if disease_id in train_disease_ids:
                        temp_train_data.append(entry)
                    else:
                        test_data.append(entry)

        # Create a new list to store the modified train data
        train_data = []

        # For each entry in the original train data...
        for entry in temp_train_data:
            # For each paragraph...
            for para in entry["paragraphs"]:
                # For each question and answer set...
                for qa in para["qas"]:
                    # For each answer...
                    for idx, answer in enumerate(qa["answers"]):
                        # Create a new question and answer set for this answer
                        new_qa = {
                            "id": qa["id"] + f"-{idx}",
                            "question": qa["question"],
                            "answers": [answer],  # Wrap the answer in a list
                            "is_impossible": qa["is_impossible"],
                        }
                        
                        # Create a new entry for this question and answer set
                        new_entry = {
                            "title": entry["title"],
                            "paragraphs": [{
                                "context": para["context"],
                                "qas": [new_qa],  # Wrap the question and answer set in a list
                            }],
                        }
                        
                        # Add the new entry to the new train data
                        train_data.append(new_entry)

        
        self._logger.info(f"Loaded all {len(temp_train_data) + len(test_data)} QAs items")

        # Analyze the lengths of all the answers in the training data
        answer_lengths = []
        for entry in train_data:
            for para in entry['paragraphs']:
                for qa in para['qas']:
                    for ans in qa['answers']:
                        answer_lengths.append(len(ans['text']))

        # Decide optimal max_answer_length based on the distribution
        # We set it to the 95th percentile length
        max_answer_length = int(np.percentile(answer_lengths, 95))


        tokenizer = AutoTokenizer.from_pretrained("dmis-lab/biobert-v1.1")
        model = AutoModelForQuestionAnswering.from_pretrained("dmis-lab/biobert-v1.1")

        train_df = self._create_df(train_data)
        test_df = self._create_df(test_data)

        train_dataset = Dataset.from_pandas(train_df)
        test_dataset = Dataset.from_pandas(test_df)

        raw_datasets = DatasetDict({
            'train': train_dataset,
            'validation': test_dataset
        })

        max_length = 384
        stride = 128

        def preprocess_training_examples(examples):
            questions = [q.strip() for q in examples["question"]]
            inputs = tokenizer(
                questions,
                examples["context"],
                max_length=max_length,
                truncation="only_second",
                stride=stride,
                return_overflowing_tokens=True,
                return_offsets_mapping=True,
                padding="max_length",
            )

            offset_mapping = inputs.pop("offset_mapping")
            sample_map = inputs.pop("overflow_to_sample_mapping")
            answers = examples["answers"]
            start_positions = []
            end_positions = []

            for i, offset in enumerate(offset_mapping):
                sample_idx = sample_map[i]
                answer = answers[sample_idx]
                start_char = answer["answer_start"][0]
                end_char = answer["answer_start"][0] + len(answer["text"][0])
                sequence_ids = inputs.sequence_ids(i)

                # Find the start and end of the context
                idx = 0
                while sequence_ids[idx] != 1:
                    idx += 1
                context_start = idx
                while sequence_ids[idx] == 1:
                    idx += 1
                context_end = idx - 1

                # If the answer is not fully inside the context, label is (0, 0)
                if offset[context_start][0] > start_char or offset[context_end][1] < end_char:
                    start_positions.append(0)
                    end_positions.append(0)
                else:
                    # Otherwise it's the start and end token positions
                    idx = context_start
                    while idx <= context_end and offset[idx][0] <= start_char:
                        idx += 1
                    start_positions.append(idx - 1)

                    idx = context_end
                    while idx >= context_start and offset[idx][1] >= end_char:
                        idx -= 1
                    end_positions.append(idx + 1)

            inputs["start_positions"] = start_positions
            inputs["end_positions"] = end_positions
            return inputs
        
        def preprocess_validation_examples(examples):
            questions = [q.strip() for q in examples["question"]]
            inputs = tokenizer(
                questions,
                examples["context"],
                max_length=max_length,
                truncation="only_second",
                stride=stride,
                return_overflowing_tokens=True,
                return_offsets_mapping=True,
                padding="max_length",
            )

            sample_map = inputs.pop("overflow_to_sample_mapping")
            example_ids = []

            for i in range(len(inputs["input_ids"])):
                sample_idx = sample_map[i]
                example_ids.append(examples["id"][sample_idx])

                sequence_ids = inputs.sequence_ids(i)
                offset = inputs["offset_mapping"][i]
                inputs["offset_mapping"][i] = [
                    o if sequence_ids[k] == 1 else None for k, o in enumerate(offset)
                ]

            inputs["example_id"] = example_ids
            return inputs

        train_dataset = raw_datasets["train"].map(
            preprocess_training_examples,
            batched=True,
            remove_columns=raw_datasets["train"].column_names,
        )

        validation_dataset = raw_datasets["validation"].map(
            preprocess_validation_examples,
            batched=True,
            remove_columns=raw_datasets["validation"].column_names,
        )
        
        # Set up the training arguments and train the model
        training_args = TrainingArguments(
            output_dir="./results_qa",
            logging_dir='./logs_qa',
            per_device_train_batch_size=21,
            per_device_eval_batch_size=9,
            gradient_accumulation_steps=3,
            evaluation_strategy="no",
            save_strategy="epoch",
            learning_rate=2e-5,
            num_train_epochs=25,
            weight_decay=0.01,
            fp16=True,
            log_level="info",
        )

        trainer = Trainer(
            model=model,
            args=training_args,
            train_dataset=train_dataset,
            eval_dataset=test_dataset,
            tokenizer=tokenizer,
        )
        
        torch.cuda.empty_cache()

        trainer.train()

        # Evaluate the model
        predictions, _, _ = trainer.predict(validation_dataset)
        start_logits, end_logits = predictions
        metrics = self._compute_metrics(start_logits, end_logits, validation_dataset, raw_datasets["validation"], max_answer_length)

        print("QA Evaluation metrics:")
        print(f"max_answer_length: {max_answer_length}\n")
        print(metrics)

        # Save the trained question-answering model
        os.makedirs(self._trained_save_path, exist_ok=True)
        trainer.save_model(self._trained_save_path)
        tokenizer.save_pretrained(self._trained_save_path)

        with open(self._max_answer_length_full_file_name, 'w') as f:
            json.dump({'max_answer_length': max_answer_length}, f)
    
    def _create_df(self, data):
        rows = []
        for entry in data:
            for para in entry['paragraphs']:
                context = para['context']
                for qa in para['qas']:
                    id = qa['id']
                    question = qa['question']
                    text_answers = [ans['text'] for ans in qa['answers']]
                    answer_start = [ans['answer_start'] for ans in qa['answers']]
                    answers = {'text': text_answers, 'answer_start': answer_start}
                    row = {'id': id, 'context': context, 'question': question, 'answers': answers}
                    rows.append(row)
        return pd.DataFrame(rows)
    
    def _compute_metrics(self, start_logits, end_logits, features, examples, max_answer_length):
            metric = evaluate.load("squad")

            example_to_features = collections.defaultdict(list)
            for idx, feature in enumerate(features):
                example_to_features[feature["example_id"]].append(idx)

            predicted_answers = []
            for example in tqdm(examples):
                example_id = example["id"]
                context = example["context"]
                answers = []

                # Loop through all features associated with that example
                for feature_index in example_to_features[example_id]:
                    start_logit = start_logits[feature_index]
                    end_logit = end_logits[feature_index]
                    offsets = features[feature_index]["offset_mapping"]

                    start_indexes = np.argsort(start_logit)[-1 : -n_best - 1 : -1].tolist()
                    end_indexes = np.argsort(end_logit)[-1 : -n_best - 1 : -1].tolist()
                    for start_index in start_indexes:
                        for end_index in end_indexes:
                            # Skip answers that are not fully in the context
                            if offsets[start_index] is None or offsets[end_index] is None:
                                continue
                            # Skip answers with a length that is either < 0 or > max_answer_length
                            if end_index < start_index:
                                continue

                            text = context[offsets[start_index][0] : offsets[end_index][1]]
                            if len(text) > max_answer_length:
                                continue

                            answer = {
                                "text": text,
                                "logit_score": start_logit[start_index] + end_logit[end_index],
                            }
                            answers.append(answer)

                # Select the answer with the best score
                if len(answers) > 0:
                    best_answer = max(answers, key=lambda x: x["logit_score"])
                    predicted_answers.append(
                        {"id": example_id, "prediction_text": best_answer["text"]}
                    )
                else:
                    predicted_answers.append({"id": example_id, "prediction_text": ""})

            theoretical_answers = [{"id": ex["id"], "answers": ex["answers"]} for ex in examples]
            return metric.compute(predictions=predicted_answers, references=theoretical_answers)

    async def evaluate_async(self, device: torch.device, question: str, articles_ids: List[str]):
        # Load the trained model and tokenizer
        model = AutoModelForQuestionAnswering.from_pretrained(self._trained_save_path)
        tokenizer = AutoTokenizer.from_pretrained(self._trained_save_path)

        max_answer_length = 0
        with open(self._max_answer_length_full_file_name, 'r') as f:
            data = json.load(f)
            max_answer_length = data['max_answer_length']
        
        self._logger.info(f"max_answer_length: {max_answer_length}")

        rows = []
        for idx, article_id in enumerate(articles_ids):
            try:
                resp = await self._get_article_async(article_id)
                context = resp.article.text
                row = {'id': article_id, 'context': context, 'question': question}
                rows.append(row)
            except Exception as e:
                continue
        
        df = pd.DataFrame(rows)
        dataset = Dataset.from_pandas(df)

        max_length = 384
        stride = 128

        def preprocess_examples(examples):
            questions = [q.strip() for q in examples["question"]]
            inputs = tokenizer(
                questions,
                examples["context"],
                max_length=max_length,
                truncation="only_second",
                stride=stride,
                return_overflowing_tokens=True,
                return_offsets_mapping=True,
                padding="max_length",
            )

            sample_map = inputs.pop("overflow_to_sample_mapping")
            example_ids = []

            for i in range(len(inputs["input_ids"])):
                sample_idx = sample_map[i]
                example_ids.append(examples["id"][sample_idx])

                sequence_ids = inputs.sequence_ids(i)
                offset = inputs["offset_mapping"][i]
                inputs["offset_mapping"][i] = [
                    o if sequence_ids[k] == 1 else None for k, o in enumerate(offset)
                ]

            inputs["example_id"] = example_ids
            return inputs
        
        processed_dataset = dataset.map(
            preprocess_examples,
            batched=True,
            remove_columns=dataset.column_names,
        )

        trainer = Trainer(
            model=model,
            tokenizer=tokenizer,
        )
        
        torch.cuda.empty_cache()

        predictions, _, _ = trainer.predict(processed_dataset)
        start_logits, end_logits = predictions

        example_to_features = collections.defaultdict(list)
        for idx, feature in enumerate(processed_dataset):
            example_to_features[feature["example_id"]].append(idx)

        predicted_answers = {}
        for example in tqdm(dataset):
            example_id = example["id"]
            context = example["context"]
            answers = []

            # Loop through all features associated with that example
            for feature_index in example_to_features[example_id]:
                start_logit = start_logits[feature_index]
                end_logit = end_logits[feature_index]
                offsets = processed_dataset[feature_index]["offset_mapping"]

                start_indexes = np.argsort(start_logit)[-1 : -n_best - 1 : -1].tolist()
                end_indexes = np.argsort(end_logit)[-1 : -n_best - 1 : -1].tolist()
                for start_index in start_indexes:
                    for end_index in end_indexes:
                        # Skip answers that are not fully in the context
                        if offsets[start_index] is None or offsets[end_index] is None:
                            continue
                        # Skip answers with a length that is either < 0 or > max_answer_length
                        if end_index < start_index:
                                continue

                        text = context[offsets[start_index][0] : offsets[end_index][1]]
                        if len(text) > max_answer_length:
                            continue

                        answer = {
                            "text": text,
                            "logit_score": start_logit[start_index] + end_logit[end_index],
                        }
                        answers.append(answer)
            
            max_answers = 10
            # Select the top k answers with the best scores
            if len(answers) > 0:
                answers.sort(key=lambda x: x["logit_score"], reverse=True)
                top_k_answers = answers[:max_answers]
                for answer in top_k_answers:
                    text = answer["text"].lower()
                    score = answer["logit_score"]

                    if score < 1.0:
                        continue

                    if text in predicted_answers:
                        if score > predicted_answers[text]['score']:
                            predicted_answers[text]['score'] = score
                        predicted_answers[text]['article_ids'].add(example_id)
                    else:
                        predicted_answers[text] = {'score': score, 'article_ids': set([example_id])}

        
        if predicted_answers:
            # Find the maximum score among all answers
            max_value = max(val['score'] for val in predicted_answers.values())

            # Calculate the threshold as a percentage of the max_value
            threshold = max_value * 0.6

            # Create a new dictionary that includes entries with values greater than or equal to the threshold
            filtered_predicted_answers = {key: val for key, val in predicted_answers.items() if val['score'] >= threshold}

            # Sort the dictionary items by score, so that higher score items appear first
            sorted_filtered_predicted_answers = {k: v for k, v in sorted(filtered_predicted_answers.items(), key=lambda item: item[1]['score'], reverse=True)}

            # Create a new dictionary to store the final answers
            final_filtered_answers = {}

            for key, value in sorted_filtered_predicted_answers.items():
                should_add = True
                for existing_key, existing_value in final_filtered_answers.items():
                    # Check if the existing answer string contains the new answer string and has higher score
                    if existing_key.find(key) != -1 and existing_value['score'] > value['score']:
                        should_add = False
                        break
                    # Check if the new answer string contains the existing answer string and has higher score
                    elif key.find(existing_key) != -1 and value['score'] > existing_value['score']:
                        del final_filtered_answers[existing_key]
                        break
                if should_add:
                    final_filtered_answers[key] = value

            return final_filtered_answers
        else:
            return predicted_answers

    @retry(
    stop=stop_after_attempt(3), 
    wait=wait_fixed(2)
    )
    async def _get_article_async(self, id: str):
        return await self._graph_client.article(id)