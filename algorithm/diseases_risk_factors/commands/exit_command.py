import sys

from thesis_diseases_risk_factors.commands._command import _Command


class ExitCommand(_Command):
    def __init__(self) -> None:
        super().__init__()

    @property
    def name(self) -> str:
        return "Exit"

    @property
    def description(self) -> str:
        return "Exit the application"
    
    async def execute_async(self, options: dict) -> None:
        sys.exit()