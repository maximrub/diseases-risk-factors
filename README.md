# Automatic Extraction of Disease Risk Factors from Medical Publications

This repository contains the implementation of our novel approach to automating the identification of disease risk factors from medical literature. Leveraging advances in machine learning and natural language processing, and specifically utilizing pre-trained models in the biomedical domain, this project aims to develop a comprehensive pipeline for extracting risk factors from a vast array of medical articles.

## Overview

Our system performs the following key tasks:

1. **Article Retrieval**: Fetches relevant medical abstracts from extensive databases like PubMed.
2. **Classification**: Uses a fine-tuned binary classifier to identify articles that discuss risk factors.
3. **Extraction**: Applies a question-answering model to precisely extract textual spans containing the risk factors.

## Trained Models and Datasets

The trained models and datasets used in this project can be found at [https://huggingface.co/diseases-risk-factors](https://huggingface.co/diseases-risk-factors).

## Repository Structure

- `algorithm`: Contains the core algorithm for disease risk factor extraction.
- `graphql-server`: Backend server using GraphQL.
- `management-ui`: Frontend user interface for managing the system.
- `LICENSE`: License information for the repository.
- `README.md`: This file.
