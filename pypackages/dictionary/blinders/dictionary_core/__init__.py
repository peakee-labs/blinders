import os


def handle_dictionary():
    """Example of calling a function from another module."""
    print("called handle dictionary", os.getcwd())
    return "handle dictionary"


__all__ = ["handle_dictionary"]
