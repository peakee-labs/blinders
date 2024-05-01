from typing import Any, Dict

from pymongo import MongoClient


def init_mongo_client(url: str) -> MongoClient[Dict[str, Any]]:
    client: MongoClient[Dict[str, Any]] = MongoClient(host=url)
    client.server_info()
    return client
