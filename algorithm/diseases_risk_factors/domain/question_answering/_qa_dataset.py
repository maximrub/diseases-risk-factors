import torch
from torch.utils.data import Dataset


class _QADataset(Dataset):
    def __init__(self, data, tokenizer):
        self.data = data
        self.tokenizer = tokenizer

    def __len__(self):
        return len(self.data)

    def __getitem__(self, idx):
        entry = self.data[idx]
        title = entry["title"]
        context = entry["paragraphs"][0]["context"]
        qas = entry["paragraphs"][0]["qas"]

        questions = []
        answers = []
        start_positions = []
        end_positions = []

        for qa in qas:
            question = qa["question"]
            answer = qa["answers"][0]  # Taking only the first answer
            start = answer["answer_start"]
            end = answer["answer_start"] + len(answer["text"]) - 1

            questions.append(question)
            answers.append(answer["text"])
            start_positions.append(start)
            end_positions.append(end)

        inputs = self.tokenizer.batch_encode_plus(
            [(question, context) for question in questions],
            truncation=True,
            padding="max_length",
            max_length=512,
            return_tensors="pt"
        )

        return {
            "input_ids": inputs["input_ids"].squeeze(),
            "attention_mask": inputs["attention_mask"].squeeze(),
            "token_type_ids": inputs["token_type_ids"].squeeze(),
            "start_positions": torch.tensor(start_positions),
            "end_positions": torch.tensor(end_positions),
        }


