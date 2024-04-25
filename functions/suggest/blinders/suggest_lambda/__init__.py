import json
import os
from typing import Any, Dict

import botocore.session

from blinders.pysuggest import explain_text_in_sentence_by_gpt_v2
from blinders.pytransport.aws import LambdaTransport
from blinders.pytransport.requests import TransportRequest, type_collect_event

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


def lambda_handler(event: Dict[str, Any], context):
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
