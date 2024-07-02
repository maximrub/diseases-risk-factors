from dependency_injector import containers, providers

from thesis_diseases_risk_factors.application.models_service import ModelsService


class ApplicationContainer(containers.DeclarativeContainer):
    config = providers.Configuration()
    utils = providers.Dependency()
    http_client = providers.Dependency()
    graph_client = providers.Dependency()

    binary_classification_model = providers.Dependency()
    question_answering_model = providers.Dependency()

    models_service = providers.Factory(
        ModelsService,
        binary_classification_model,
        question_answering_model,
        graph_client=graph_client,
        http_client=http_client,
        trained_save_path=config.models.rootTrainedSavePath,
        trained_bin_path=config.models.rootTrainedBinPath,
        utils=utils
    )