class CommandOption:
    def __init__(self,
        name,
        description,
        default_value):
        self._name = name
        self._description = description
        self._default_value = default_value
    
    @property
    def name(self) -> str:
        return self._name
    
    @property
    def description(self) -> str:
        return self._description
    
    @property
    def default_value(self):
        return self._default_value