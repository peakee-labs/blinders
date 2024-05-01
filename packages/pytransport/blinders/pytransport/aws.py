import json
from typing import Any, Dict

from blinders.pytransport import ITransport


class LambdaTransport(ITransport):
    def __init__(self, client):
        self.client = client

    def request(self, id: str, payload: bytes) -> bytes:
        print("lambda transport: request to", id)
        response: Dict[str, Any] = self.client.invoke(
            FunctionName=id,
            Payload=payload,
            InvocationType="RequestResponse",
        )
        if response.get("FunctionError"):
            res = json.loads(response["Payload"].read().decode("utf-8"))
            raise Exception(json.dumps(res))

        return response["Payload"].read()

    def push(self, id: str, payload: bytes):
        print("lambda transport: push to", id)
        self.client.invoke(
            FunctionName=id,
            Payload=payload,
            InvocationType="Event",
        )
