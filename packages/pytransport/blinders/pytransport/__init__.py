class ITransport:
    def request(self, id: str, payload: bytes) -> bytes:
        raise Exception("this method is not implemented")

    def push(self, id: str, payload: bytes) -> None:
        raise Exception("this method is not implemented")
