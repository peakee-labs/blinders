import typing

type_collect_event = "COLLECT_EVENT"


class TransportRequest:
    Type: str
    Data: dict[str, typing.Any]

    def __init__(self, type: str, data: dict[str, typing.Any]) -> None:
        self.Type = type
        self.Data = data

    def json(self) -> dict[str, typing.Any]:
        return {"type": self.Type, "data": self.Data}
