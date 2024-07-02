import logging
import uuid
from typing import List
import aiohttp
import aiofiles
import torch
import sys
import json
import shutil
import os

from thesis_diseases_risk_factors.graph.client import Client
from thesis_diseases_risk_factors.graph.client.input_types import RiskFactorInput
from thesis_diseases_risk_factors.domain.binary_classification.binary_classification_model import BinaryClassificationModel
from thesis_diseases_risk_factors.domain.question_answering.question_answering_model import QuestionAnsweringModel
from thesis_diseases_risk_factors.domain.shared_kernel.utils import Utils


class ModelsService:
    def __init__(self,
        binary_classification_model: BinaryClassificationModel,
        question_answering_model: QuestionAnsweringModel,
        graph_client: Client,
        http_client: aiohttp.ClientSession,
        trained_save_path: str,
        trained_bin_path: str,
        utils: Utils) -> None:
        self._logger = logging.getLogger(__name__)
        self._binary_classification_model = binary_classification_model
        self._question_answering_model = question_answering_model
        self._graph_client = graph_client
        self._http_client = http_client
        self._trained_save_path = trained_save_path
        self._trained_bin_path = trained_bin_path
        self._utils = utils
    
    async def train_async(self):
        # Set up the device
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self._logger.info("device used: [%s]", device.type)

        await self._binary_classification_model.train_async(device=device)
        await self._question_answering_model.train_async(device=device)

        output_filename = f"{str(uuid.uuid4())}"
        output_filename_with_ext = f"{output_filename}.tar.gz"
        self._utils.gzip_folder(self._trained_save_path, output_filename)

        if not os.path.exists(self._trained_bin_path):
            os.mkdir(self._trained_bin_path)

        current_dir = os.getcwd()
        shutil.move(f"{current_dir}/{output_filename_with_ext}", f"{current_dir}/{self._trained_bin_path}/{output_filename_with_ext}")
        
    async def evaluate_async(self):
        # Set up the device
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self._logger.info("device used: [%s]", device.type)


        diseases_resp = await self._graph_client.list_diseases()
        for disease in diseases_resp.diseases:
            disease_name = disease.names[0]
            self._logger.info(f"Learning risk factors for {disease.id} :: {disease_name}")
            
            # Build the search term string
            search_term = f'"{disease_name}"[Title/Abstract/MeSH Terms] AND "Risk Factors"[Title/Abstract/MeSH Terms]'
            articles_resp = await self._graph_client.search_articles(search_term, 1000)
            articles_ids = articles_resp.search_articles
            if articles_ids:
                risk_factors_ids = await self._binary_classification_model.evaluate_async(device=device, articles_ids=articles_ids) 
                if risk_factors_ids:
                    question = f"What are the risk factors of {disease_name}?"
                    answers_dict = await self._question_answering_model.evaluate_async(device=device, question=question, articles_ids=risk_factors_ids)

                    if answers_dict:
                        self._logger.info(f"Upload risk factors for {disease.id} :: {disease_name}")
                        risk_factors = []
                        for k, v in answers_dict.items():
                            risk_factors.append(RiskFactorInput(text=k, score=v["score"], articlesIds=list(v["article_ids"])))
                        await self._graph_client.update_risk_factors(disease_id=disease.id, risk_factors=risk_factors)