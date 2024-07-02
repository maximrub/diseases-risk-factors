from aiohttp import http
from dependency_injector import containers, providers

from thesis_diseases_risk_factors.application import ApplicationContainer
from thesis_diseases_risk_factors.domain.shared_kernel import SharedKernelContainer
from thesis_diseases_risk_factors.infrastructure.resources import ResourcesContainer
from thesis_diseases_risk_factors.domain.binary_classification import BinaryClassificationContainer
from thesis_diseases_risk_factors.domain.question_answering import QuestionAnsweringContainer
from thesis_diseases_risk_factors.commands import CommandsContainer
from thesis_diseases_risk_factors.menu_invoker import MenuInvoker

class RootContainer(containers.DeclarativeContainer):
    config = providers.Configuration()

    resources_package = providers.Container(
        ResourcesContainer,
        config=config.resources
    )

    shared_kernel_package = providers.Container(
        SharedKernelContainer,
        config=config.shared_kernel.domain
    )

    binary_classification_domain_package = providers.Container(
        BinaryClassificationContainer,
        config=config.models,
        utils=shared_kernel_package.utils,
        graph_client=resources_package.graph_client,
        http_client=resources_package.http_client
    )

    question_answering_classification_domain_package = providers.Container(
        QuestionAnsweringContainer,
        config=config.models,
        utils=shared_kernel_package.utils,
        graph_client=resources_package.graph_client,
        http_client=resources_package.http_client
    )

    application_package = providers.Container(
        ApplicationContainer,
        config=config,
        utils=shared_kernel_package.utils,
        http_client=resources_package.http_client,
        graph_client=resources_package.graph_client,
        binary_classification_model = binary_classification_domain_package.binary_classification_model,
        question_answering_model = question_answering_classification_domain_package.question_answering_model
    )

    commands_package = providers.Container(
        CommandsContainer,
        models_service=application_package.models_service,
    )

    menu_invoker = providers.Factory(
        MenuInvoker,
        commands=commands_package.commands
    )