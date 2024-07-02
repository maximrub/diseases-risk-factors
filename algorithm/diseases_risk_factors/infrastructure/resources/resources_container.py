import logging.config
from dependency_injector import containers, providers

from thesis_diseases_risk_factors.infrastructure.resources.http_client_resource import HttpClientResource
from thesis_diseases_risk_factors.infrastructure.resources.graph_client_resource import GraphClientResource


class ResourcesContainer(containers.DeclarativeContainer):
    config = providers.Configuration()

    logging = providers.Resource(
        logging.config.dictConfig,
        config=config.logging,
    )

    http_client = providers.Resource(
        HttpClientResource.init_async,
        config=config.httpClient
    )

    graph_client = providers.Resource(
        GraphClientResource.init,
        config=config.graph
    )