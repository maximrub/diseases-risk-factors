from dependency_injector import containers, providers

from thesis_diseases_risk_factors.domain.shared_kernel.utils import Utils 


class SharedKernelContainer(containers.DeclarativeContainer):
    config = providers.Configuration()

    utils = providers.Factory(
        Utils
    )
