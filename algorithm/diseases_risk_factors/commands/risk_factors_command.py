from typing import Dict

from thesis_diseases_risk_factors.commands._command import _Command
from thesis_diseases_risk_factors.commands.entities import CommandOption
from thesis_diseases_risk_factors.application.models_service import ModelsService


class RiskFactorsCommand(_Command):
    def __init__(self,
        models_service: ModelsService) -> None:
        super().__init__()
        self._models_service = models_service

    @property
    def name(self) -> str:
        return "Write Risk Factors"

    @property
    def description(self) -> str:
        return "Write Risk Factors"
    
    async def execute_async(self, options: dict) -> None:
        await self._models_service.evaluate_async()