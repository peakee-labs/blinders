from datetime import datetime
from typing import Any

from dotenv import load_dotenv
from langchain_core.output_parsers import JsonOutputParser
from langchain_openai import ChatOpenAI

from blinders.pysuggest.prompts import translate_words_in_sentence


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


def main():
    load_dotenv()

    text = "developing an extension"
    sentence = "If you are developing an extension that talks with some API you\
          probably are using different keys for testing and production"
    result = explain_text_in_sentence_by_gpt(text, sentence)

    # {
    #     'translate': 'phát triển một tiện ích mở rộng',
    #     'grammar_analysis': {
    #          'tense': {'type': 'Present continuous', 'identifier': 'are developing'}
    #     },
    #     'expand_words': ['programming', 'coding', 'implementing']
    # }

    print(result)
