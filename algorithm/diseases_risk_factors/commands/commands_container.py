from dependency_injector import containers, providers


from thesis_diseases_risk_factors.commands.train_models_command import TrainModelsCommand
from thesis_diseases_risk_factors.commands.risk_factors_command import RiskFactorsCommand
from thesis_diseases_risk_factors.commands.exit_command import ExitCommand


class CommandsContainer(containers.DeclarativeContainer):
    models_service = providers.Dependency()


    commands = providers.List(
            providers.Factory(TrainModelsCommand, models_service=models_service),
            providers.Factory(RiskFactorsCommand, models_service=models_service),
            providers.Factory(ExitCommand),
    )