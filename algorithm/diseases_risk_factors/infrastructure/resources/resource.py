from abc import ABCMeta, abstractmethod
from typing import TypeVar, Generic


T = TypeVar('T')


class Resource(Generic[T], metaclass=ABCMeta):

    @classmethod
    def init(cls, config) -> T:
        resource = cls._create(config)
        yield resource
        cls._shutdown(resource)


    @classmethod
    @abstractmethod
    def _create(cls, config) -> T:
        pass

    @classmethod
    @abstractmethod
    def _shutdown(cls, resource: T) -> None:
        pass