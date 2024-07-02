import logging
from typing import List

from thesis_diseases_risk_factors import commands
from thesis_diseases_risk_factors.commands import _Command


class MenuInvoker:
    def __init__(self,
        commands: List[_Command]) -> None:
        self._logger = logging.getLogger(__name__)
        self._commands = {index: command for index, command in enumerate(commands, 1)}

    async def execute_async(self, command_index) -> None:
        if command_index in self._commands:
            options = {}
            for option_name, option in self._commands[command_index].command_options.items():
                options[option.name] = input(f"{option.description} (Enter for None): ")
            result = await self._commands[command_index].execute_async(options)
            if result is not None:
                print(result)
        else:
            print(f"{command_index } isn't a valid command index")
    
    def show(self) -> None:
        for index, command in self._commands.items():
            print(f"{index}: {command.name}")