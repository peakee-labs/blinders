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
