import aiohttp
from typing import get_type_hints

from thesis_diseases_risk_factors.infrastructure.resources.async_resource import AsyncResource


class HttpClientResource(AsyncResource):
    @classmethod
    async def _create_async(cls, config) -> aiohttp.ClientSession:
        client: aiohttp.ClientSession
        if config["forceClose"]:
            connector = aiohttp.TCPConnector(force_close=True)
            client = aiohttp.ClientSession(connector=connector)
        else:
            client = aiohttp.ClientSession()
        
        return client

    @classmethod
    async def _shutdown_async(cls, resource: aiohttp.ClientSession) -> None:
        await resource.close()