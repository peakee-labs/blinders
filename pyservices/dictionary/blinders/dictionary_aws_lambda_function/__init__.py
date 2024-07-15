from blinders.dictionary_core import handle_dictionary


def call_handle_dictionary():
    """Example of calling a function from another module."""
    return "called " + handle_dictionary()


def lambda_handler(event, context):
    """Example of calling a function from another module."""
    print("dictionary on event", event["resource"])
    return {"statusCode": 200, "body": call_handle_dictionary()}
