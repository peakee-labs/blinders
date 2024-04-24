from blinders.pytransport import ITransport


class LambdaTransport(ITransport):
    def __init__(self, client):
        self.client = client

    def request(self, id: str, payload: bytes) -> bytes:
        print("lambda transport: request to", id)
        response = self.client.invoke(
            FunctionName=id,
            Payload=payload,
            InvocationType="RequestResponse",
        )
        if response["FunctionError"] != "":
            raise Exception(
                "lambda transport: cannot invoke, err",
                response["FunctionError"],
            )

        return response["Payload"].read()

    def push(self, id: str, payload: bytes):
        print("lambda transport: push to", id)
        self.client.invoke(
            FunctionName=id,
            Payload=payload,
            InvocationType="Event",
        )
