class ArgumentError(Exception):
    def __init__(self, name, value, *args):
        """
        Raise when an argument value is invalid
        """

        self.name = name
        self.value = value
        super().__init__(f"Value: [{value}] for argument [{name}] is invalid", *args)