from dependency_injector import containers, providers

from thesis_diseases_risk_factors.domain.question_answering.question_answering_model import QuestionAnsweringModel


class QuestionAnsweringContainer(containers.DeclarativeContainer):
    config = providers.Configuration()
    utils = providers.Dependency()
    graph_client = providers.Dependency()
    http_client = providers.Dependency()

    question_answering_model = providers.Factory(
        QuestionAnsweringModel,
        graph_client=graph_client,
        http_client=http_client,
        trained_save_path=config.questionAnsweringTrainedSavePath
    )
