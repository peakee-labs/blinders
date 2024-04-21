import json
from datetime import datetime
from typing import Any

from dotenv import load_dotenv
from langchain_core.output_parsers import JsonOutputParser
from langchain_openai import ChatOpenAI

from blinders.pysuggest.prompts import (
    explain_words_in_sentence,
    translate_words_in_sentence,
)


def explain_text_in_sentence_by_gpt(text: str, sentence: str):
    print("start querying explain_text_in_sentence_by_gpt")
    model = ChatOpenAI(model="gpt-3.5-turbo")  # gpt-4 takes too long
    output_parser = JsonOutputParser()
    chain = translate_words_in_sentence | model | output_parser

    start = datetime.now()
    result: dict[str, Any] = chain.invoke({"text": text, "sentence": sentence})
    end = datetime.now()
    result["duration_in_seconds"] = (end - start).total_seconds()
    print("end explain_text_in_sentence_by_gpt")

    return result


def explain_text_in_sentence_by_gpt_v2(text: str, sentence: str):
    model = ChatOpenAI(model="gpt-3.5-turbo")  # gpt-4 takes too long
    output_parser = JsonOutputParser()
    chain = explain_words_in_sentence | model | output_parser

    start = datetime.now()
    result: dict[str, Any] = chain.invoke({"text": text, "sentence": sentence})
    end = datetime.now()
    result["duration_in_seconds"] = (end - start).total_seconds()

    return result


def main():
    load_dotenv()

    text = "developing an extension"
    sentence = "If you are developing an extension that talks with some API you\
          probably are using different keys for testing and production"

    # {
    #     'translate': 'phát triển một tiện ích mở rộng',
    #     'grammar_analysis': {
    #          'tense': {'type': 'Present continuous', 'identifier': 'are developing'}
    #     },
    #     'expand_words': ['programming', 'coding', 'implementing']
    # }

    text = "open-source project funding"
    sentence = "I know that open-source project funding has been on people's minds lately, \
        as we've watched the Internet's reaction\
              to the bizarre open-sourcing of the V language unfold"

    result = explain_text_in_sentence_by_gpt_v2(text, sentence)

    print(json.dumps(result, indent=2))

    #   "translate": "kho\u1ea3n t\u00e0i tr\u1ee3 d\u1ef1 \u00e1n m\u00e3 ngu\u1ed3n m\u1edf",
    #   "IPA": "\u02c8o\u028ap\u0259n-s\u0254\u02d0rs \u02c8pr\u0252d\u0292ekt..."
    #   "grammar_analysis": {
    #     "tense": {
    #       "type": "present perfect continuous",
    #       "identifier": "has been"
    #     },
    #     "structure": {
    #       "type": "complex sentence",
    #       "structure": "I know that + S + has been + V_ed",
    #       "for": "showing a continuous action happening over a period of time"
    #     }
    #   },
    #   "key_words": [
    #     "Internet",
    #     "reaction"
    #   ],
    #   "expand_words": [
    #     "bizarre",
    #     "V language",
    #     "unfold"
    #   ],
    #   "duration_in_seconds": 2.858413
    # }
