from abc import ABCMeta, abstractmethod
from typing import TypeVar, Generic


T = TypeVar('T')


class AsyncResource(Generic[T], metaclass=ABCMeta):

    @classmethod
    async def init_async(cls, config) -> T:
        resource = await cls._create_async(config)
        yield resource
        await cls._shutdown_async(resource)

    @classmethod
    @abstractmethod
    async def _create_async(cls, config) -> T:
        pass

    @classmethod
    @abstractmethod
    async def _shutdown_async(cls, resource: T) -> None:
        pass