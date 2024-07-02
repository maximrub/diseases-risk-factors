from abc import ABCMeta, abstractmethod
import logging
from typing import Dict

from thesis_diseases_risk_factors.commands.entities import CommandOption


class _Command(metaclass=ABCMeta):
    def __init__(self) -> None:
        self._logger = logging.getLogger(type(self).__name__)
    
    @property
    @abstractmethod
    def name(self) -> str:
        pass

    @property
    @abstractmethod
    def description(self) -> str:
        pass

    @property
    def command_options(self) -> Dict[str, CommandOption]:
        return {}

    @abstractmethod
    async def execute_async(self, options):
        pass
