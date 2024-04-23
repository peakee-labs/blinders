class ITransport:
    def Request(self, id: str, payload: bytes) -> bytes:
        raise Exception("this method is not implemented")

    def Push(self, id: str, payload: bytes) -> None:
        raise Exception("this method is not implemented")
