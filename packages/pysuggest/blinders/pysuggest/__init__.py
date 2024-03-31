from datetime import datetime

from dotenv import load_dotenv
from langchain_core.output_parsers import JsonOutputParser
from langchain_openai import ChatOpenAI

from blinders.pysuggest.prompts import translate_words_in_sentence


def main():
    load_dotenv()

    model = ChatOpenAI(model="gpt-3.5-turbo")  # gpt-4 takes too long
    output_parser = JsonOutputParser()
    chain = translate_words_in_sentence | model | output_parser

    start = datetime.now()
    result = chain.invoke(
        {
            "text": "developing an extension",
            "sentence": "If you are developing an extension that talks with some API you probably\
                  are using different keys for testing and production",
        }
    )
    end = datetime.now()

    # {
    #     'translate': 'phát triển một tiện ích mở rộng',
    #     'grammar_analysis': {
    #          'tense': {'type': 'Present continuous', 'identifier': 'are developing'}
    #     },
    #     'expand_words': ['programming', 'coding', 'implementing']
    # }

    print(result)
    print((end - start).total_seconds(), "s")

    # 'explain': 'very short explain the meaning of text in the sentence',\
