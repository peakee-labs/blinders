from langchain_core.prompts import ChatPromptTemplate

translate_words_in_sentence = ChatPromptTemplate.from_template(
    """You are an English teacher. \
    The student want to understand the word "{text}" in "{sentence}". \
    Return json with the following format: \
    {{\
        "translate": word to Vietnamese,\
        "grammar_analysis": {{"tense": {{"type": type of tense, "identifier": identifier}}}},\
        "expand_words": give 3 random words in this major\
    }}\
    Make it short, clear and easy to understand.\
    """
)

explain_words_in_sentence = ChatPromptTemplate.from_template(
    """You are an English professor. \
    The Vietnamese want to learn the words "{text}" in "{sentence}". \
    Return json with the following format: \
    {{\
        "translate": word to Vietnamese,\
        "IPA": IPA English pronunciation,\
        "grammar_analysis": {{\
            "tense": {{"type": type of tense of the whole sentence, "identifier": identifier}},\
            "structure": {{"type": structure type of the whole sentence,\
                "structure": show the structure of the sentence \
                    as form example 'I know that + S + has been + V_ed +', \
                "for": how to structure is used for
            }} \
        }},
        "key_words": get 1 or 2 or 3 main words in the sentence not in the text\
              (it can be noun/verb/adjective)\
        "expand_words": give 3 words might be relevant but not in the sentence\
    }}\
    Make it short, clear and easy to understand.\
    """
)
