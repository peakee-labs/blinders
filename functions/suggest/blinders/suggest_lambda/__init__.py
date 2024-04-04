import json
from typing import Any, Dict

from blinders.pysuggest import explain_text_in_sentence_by_gpt

request_types = ["explain-text-in-sentence"]
models = ["gpt"]


def lambda_handler(event: Dict[str, Any], context):
    """Example of calling a function from another module."""
    queries: Dict[str, str] = event["queryStringParameters"]
    print("handle suggest with payload", queries)

    if queries is None:
        return {"statusCode": 400, "body": "require queries"}

    tp = queries.get("type")
    if tp not in request_types:
        return {
            "statusCode": 400,
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
                    "body": "text and sentence are required for gpt",
                }
            suggest = explain_text_in_sentence_by_gpt(text, sentence)
            print("suggest", suggest)
            return {"statusCode": 200, "body": json.dumps(suggest)}
