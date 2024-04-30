import json
import os
from typing import Any, Dict


import botocore.session

from blinders.pysuggest import explain_text_in_sentence_by_gpt_v2
from blinders.pytransport.aws import LambdaTransport
from blinders.pytransport.requests import TransportRequest, type_collect_event
from blinders.pydb import MongoManager
from blinders.pydb.utils.mongo import init_mongo_client

request_types = ["explain-text-in-sentence"]
models = ["gpt"]

default_headers = {
    "Access-Control-Allow-Origin": "*",
    "Content-Type": "application/json",
}

session = botocore.session.get_session()
client = session.create_client("lambda")
transport = LambdaTransport(client)
collect = "COLLECT"
authenticate = "AUTHENTICATE"
consumeMap = {
    collect: os.getenv("COLLECTING_PUSH_FUNCTION_NAME", ""),
    authenticate: os.getenv("AUTHENTICATE_FUNCTION_NAME"),
}

db_url = "mongodb://{}:{}@{}:{}/{}".format(
    os.getenv("MONGO_USERNAME"),
    os.getenv("MONGO_PASSWORD"),
    os.getenv("MONGO_HOST"),
    os.getenv("MONGO_PORT"),
    os.getenv("MONGO_DATABASE"),
)

mongo_manager = MongoManager(
    client=init_mongo_client(db_url),
    name=os.getenv("MONGO_DATABASE", ""),
)


def lambda_handler(event: Dict[str, Any], context):
    """Example of calling a function from another module."""
    headers: Dict[str, Any] = event.get("headers", None)
    if headers is None:
        raise Exception("headers field not exsisted in event")

    auth: str = headers.get("authorization", None)
    if auth is None or auth == "":
        return {
            "statusCode": 400,
            "body": "missing authorization header",
            "headers": default_headers,
        }

    auth_request = {"token": auth}
    auth_user: Dict[str, Any]
    try:
        payload = json.dumps(auth_request).encode("utf-8")
        response = transport.request(consumeMap[authenticate], payload)
        data_str = response.decode("utf-8")
        auth_user = json.loads(data_str)
        print("pysuggest: auth user from authenticate", auth_user)

    except Exception as e:
        try:
            print("pysuggest: failed to authenticate with authenticate lambda ", e)
            return {
                "statusCode": 400,
                "headers": default_headers,
                "body": str(e),
            }

        except Exception as er:
            print("pysuggest: failed to parse error from authenticate lambda", er)
            return {
                "statusCode": 400,
                "headers": default_headers,
                "body": "failed to verify given token",
            }

    queries: Dict[str, str] = event["queryStringParameters"]
    print("handle suggest with payload", queries)

    if queries is None:
        return {"statusCode": 400, "body": "require queries", "headers": default_headers}

    tp = queries.get("type")
    if tp not in request_types:
        return {
            "statusCode": 400,
            "headers": default_headers,
            "body": "unsupported type, expect" + str(request_types),
        }

    if tp == "explain-text-in-sentence":
        model = queries.get("model") or "gpt"
        if model not in models:
            return {
                "statusCode": 400,
                "body": "unsupported type, expect" + str(models),
            }

        if model == "gpt":
            text = queries.get("text")
            sentence = queries.get("sentence")
            if text is None or sentence is None:
                return {
                    "statusCode": 400,
                    "headers": default_headers,
                    "body": "text and sentence are required for gpt",
                }
            suggest = explain_text_in_sentence_by_gpt_v2(text, sentence)

            # TODO: make struct
            suggest_event = {
                "userId": auth_user.get("ID"),
                "request": {
                    "text": text,
                    "sentence": sentence,
                },
                "response": {
                    "translate": suggest.get("translate", ""),
                    "grammarAnalysis": suggest.get("grammar_analysis", ""),
                    "expandWords": suggest.get("expand_words", ""),
                    "keyWords": suggest.get("key_words", ""),
                },
            }
            generic_event = {
                "type": "EXPLAIN",
                "event": suggest_event,
            }
            transport_event = TransportRequest(type=type_collect_event, data=generic_event)

            try:
                payload = json.dumps(transport_event.json()).encode("utf-8")
                transport.push(consumeMap[collect], payload)

            except Exception as e:
                print("pysuggest: cannot push to collecting", e)

            return {"statusCode": 200, "headers": default_headers, "body": json.dumps(suggest)}
