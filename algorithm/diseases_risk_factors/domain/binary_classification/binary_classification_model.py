import logging
import os
from typing import List
import aiohttp
import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification, TrainingArguments, Trainer
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
import numpy as np
import pandas as pd
from datasets import Dataset
from tenacity import retry, stop_after_attempt, wait_fixed

from thesis_diseases_risk_factors.domain.binary_classification._classification_dataset import _ClassificationDataset
from thesis_diseases_risk_factors.domain.binary_classification._binary_classifier import _BinaryClassifier
from thesis_diseases_risk_factors.graph.client import Client


class BinaryClassificationModel:
    def __init__(self,
        graph_client: Client,
        http_client: aiohttp.ClientSession,
        trained_save_path: str) -> None:
        self._logger = logging.getLogger(__name__)
        self._graph_client = graph_client
        self._http_client = http_client
        self._trained_save_path = trained_save_path
        self._model_name = "risk_factors_binary_classification.pth"

    async def train_async(self, device: torch.device):
        response = await self._graph_client.list_classification_items()
        items = response.classification_items

        # Separate the texts and the labels
        texts = [item.article.text for item in items]
        labels = [item.label for item in items]

        # Split the data into train and test sets
        train_texts, test_texts, train_labels, test_labels = train_test_split(texts, labels, test_size=0.2)

        tokenizer = AutoTokenizer.from_pretrained("dmis-lab/biobert-v1.1")
        model = AutoModelForSequenceClassification.from_pretrained("dmis-lab/biobert-v1.1")

        # Prepare the data
        train_encodings = tokenizer(train_texts, truncation=True, padding=True, max_length=512)
        test_encodings = tokenizer(test_texts, truncation=True, padding=True, max_length=512)

        # Create instances of your ClassificationDataset for training and testing data
        train_dataset = _ClassificationDataset(train_encodings, train_labels)
        test_dataset = _ClassificationDataset(test_encodings, test_labels)
        
        # Initialize the Trainer
        training_args = TrainingArguments(
            output_dir='./results_binary',
            logging_dir='./logs_binary',
            num_train_epochs=3,
            per_device_train_batch_size=3,
            per_device_eval_batch_size=3,
            gradient_accumulation_steps=3,
            log_level="info",
        )

        trainer = Trainer(
            model=model,
            args=training_args,
            train_dataset=train_dataset,
            eval_dataset=test_dataset,
        )

        torch.cuda.empty_cache()

        # Train the model
        trainer.train()

        # Get predictions and ground truth for test dataset
        predictions, labels, _ = trainer.predict(test_dataset)
        predictions = np.argmax(predictions, axis=1)

        # Print the classification report
        self._logger.info("binary classification report")
        print(classification_report(labels, predictions))

        # Save the trained question-answering model
        os.makedirs(self._trained_save_path, exist_ok=True)
        trainer.save_model(self._trained_save_path)
        tokenizer.save_pretrained(self._trained_save_path)
    
    async def evaluate_async(self, device: torch.device, articles_ids: List[str]) -> List[str]:
         # Load the trained model and tokenizer
        model = AutoModelForSequenceClassification.from_pretrained(self._trained_save_path)
        tokenizer = AutoTokenizer.from_pretrained(self._trained_save_path)

        trainer = Trainer(
            model=model,
            tokenizer=tokenizer,
        )

        df = pd.DataFrame(columns=['id', 'text'])

        for id in articles_ids:
            try:
                resp = await self._get_article_async(id)
            except Exception as e:
                continue
            
            new_row = pd.DataFrame({'id': [id], 'text': [resp.article.text]})
            df = pd.concat([df, new_row], ignore_index=True)
        
        dataset = Dataset.from_pandas(df)

        def tokenize_text(example):
            return tokenizer(example['text'], truncation=True, padding=True, max_length=512)

        dataset = dataset.map(tokenize_text, batched=True)

        torch.cuda.empty_cache()

        predictions = trainer.predict(dataset)
        if predictions is not None and predictions.predictions is not None:
            predicted_labels = predictions.predictions.argmax(axis=1)
            df['label'] = predicted_labels

            ids_with_label_1 = df.loc[df['label'] == 1, 'id'].tolist()
            return ids_with_label_1
        else:
            return None
    
    @retry(
    stop=stop_after_attempt(3), 
    wait=wait_fixed(2)
    )
    async def _get_article_async(self, id: str):
        return await self._graph_client.article(id)
