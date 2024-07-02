from dependency_injector import containers, providers

from thesis_diseases_risk_factors.domain.binary_classification.binary_classification_model import BinaryClassificationModel


class BinaryClassificationContainer(containers.DeclarativeContainer):
    config = providers.Configuration()
    utils = providers.Dependency()
    graph_client = providers.Dependency()
    http_client = providers.Dependency()

    binary_classification_model = providers.Factory(
        BinaryClassificationModel,
        graph_client=graph_client,
        http_client=http_client,
        trained_save_path=config.binaryClassificationTrainedSavePath
    )
