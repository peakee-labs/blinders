class ITransport:
    def Request(self, id: str, payload: str) -> str:
        raise Exception("this method is not implemented")

    def Push(self, id: str, payload: str) -> None:
        raise Exception("this method is not implemented")
