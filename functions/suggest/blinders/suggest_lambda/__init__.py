import json
import os
from typing import Any, Dict

import botocore.session

from blinders.pysuggest import explain_text_in_sentence_by_gpt_v2
from blinders.pytransport.aws import LambdaTransport
from blinders.pytransport.requests import TransportRequest, type_collect_event
from blinders.pyauth import AuthUser, AuthManager
from blinders.pyauth.aws import auth
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
consumeMap = {collect: os.getenv("COLLECTING_PUSH_FUNCTION_NAME", "")}


auth_cert: Dict[str, Any] = {}
with open("firebase.admin.json") as file:
    auth_cert = json.load(file)

auth_manager = AuthManager(auth_cert)

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


@auth(auth_manager=auth_manager, repo=mongo_manager.UsersRepo)
def lambda_handler(event: Dict[str, Any], context, auth_user: AuthUser):
    """Example of calling a function from another module."""

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
