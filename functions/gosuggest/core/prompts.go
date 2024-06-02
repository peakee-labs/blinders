package core

type ExplainPhraseInSentenceResult struct {
	Translate       string `json:"translate"`
	IPA             string `json:"IPA"`
	GrammarAnalysis struct {
		Tense struct {
			Type       string
			Identifier string
		}
		Structure struct {
			Type      string
			Structure string
			For       string
		}
	} `json:"grammarAnalysis"`
	KeyWords    []string `json:"keyWords"`
	ExpandWords []string `json:"expandWords"`
}

const ExplainPhraseInSentencePrompt = `
<|begin_of_text|>

<|start_header_id|>system<|end_header_id|> 
You are a helpful AI assistant for learning English for Vietnamese learner.  Only response to user a JSON object, the result must be short, clear. You must response the result with user's instruction, and the result must be correct with user's instruction. Don't give any other information or make any action without user's instruction.
<|eot_id|>

<|start_header_id|>user<|end_header_id|>
I want to understand:
- phrase: "%v"
- sentence: "%v"
<|eot_id|>

<|start_header_id|>system<|end_header_id|> 
Only respond JSON format: 
    {
        "translate": a string translate only 'the phrase' not 'the sentence' to Vietnamese,
        "IPA": IPA English pronunciation of phrase,
        "grammarAnalysis": {
            "tense": {"type": type of tense of the whole sentence, "identifier": how user can identify the tense},
            "structure": {"type": structure type of the whole sentence,
                "structure": show the grammar structure of the sentence as form example 'I know that + S + has been + V_ed +',
                "for": how the structure is used for
            }
        },
        "keyWords": get 1, or 2, or 3 main words in the sentence (it can be noun/verb/adjective),
        "expandWords": give 3 words might be relevant but not in the sentence
    }
<|eot_id|>

<|start_header_id|>assistant<|end_header_id|>
`
