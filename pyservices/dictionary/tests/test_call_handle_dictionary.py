from blinders.dictionary_aws_lambda_function import call_handle_dictionary


def test_call_handle_dictionary():
    assert call_handle_dictionary() == "called handle dictionary"
