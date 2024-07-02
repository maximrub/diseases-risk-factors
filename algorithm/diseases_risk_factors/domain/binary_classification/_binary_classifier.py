import torch.nn as nn


class _BinaryClassifier(nn.Module):
    """
    The binary classification model
    """
    def __init__(self, bert):
        super(_BinaryClassifier, self).__init__()
        self._bert = bert
        self._classifier = nn.Linear(self._bert.config.hidden_size, 2)

    def forward(self, input_ids, attention_mask):
        outputs = self._bert(input_ids=input_ids, attention_mask=attention_mask)
        logits = self._classifier(outputs.last_hidden_state[:, 0, :])
        return logits