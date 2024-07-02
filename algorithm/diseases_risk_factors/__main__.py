import logging
import asyncio
import sys
import os
from dependency_injector.wiring import inject, Provide, Closing

from thesis_diseases_risk_factors.root_container import RootContainer


@inject
async def run(
    menu_invoker = Provide[RootContainer.menu_invoker],
):
    logger = logging.getLogger(__name__)
    try:
        command_index = os.getenv('COMMAND_INDEX')
        if command_index is not None:
            await menu_invoker.execute_async(int(command_index))
        else:
            while True:
                menu_invoker.show()
                command_index = int(input("Select menu option: "))
                await menu_invoker.execute_async(command_index)
    except Exception as error:
        logger.exception(error, exc_info=True)
        raise
    

async def main(root_container):
    await root_container.init_resources()
    await run()
    await root_container.shutdown_resources()


if __name__ == "__main__":
    root_container = RootContainer()
    root_container.config.from_yaml(os.path.join(os.path.dirname(__file__), "config.yaml"), required=True)
    root_container.wire(modules=[sys.modules[__name__]])
    asyncio.run(main(root_container))