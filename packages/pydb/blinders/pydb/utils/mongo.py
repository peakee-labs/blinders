from pymongo import MongoClient
from typing import Dict, Any


def init_mongo_client(url: str) -> MongoClient[Dict[str, Any]]:
    client: MongoClient[Dict[str, Any]] = MongoClient(host=url)
    client.server_info()
    return client
