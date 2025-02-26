class {{.ClientName}}Exception(Exception):

    def __init__(self, message, status):
        super().__init__(message, status)
        self.message = message
        self.status = status