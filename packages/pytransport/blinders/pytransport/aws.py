import boto3

from blinders.pytransport import ITransport


class LambdaTransport(ITransport):
    def __init__(self, *args):
        self.client = boto3.client("lambda", args)

    def Request(self, id: str, payload: str) -> str:
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

        return response["Payload"]

    def Push(self, id: str, payload: str):
        print("lambda transport: push to", id)
        response = self.client.invoke(
            FunctionName=id,
            Payload=payload,
            InvocationType="Event",
        )

        if response["FunctionError"] != "":
            raise Exception(
                "lambda transport: cannot invoke, err",
                response["FunctionError"],
            )
