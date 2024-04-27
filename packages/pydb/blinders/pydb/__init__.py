import pymongo
import pymongo.database
from typing import Dict, Any


from blinders.pydb.repository.users import UsersRespo

UserColName = "users"


class MongoManager:
    Client: pymongo.MongoClient[Dict[str, Any]]
    Database: str

    UsersRepo: UsersRespo

    def __init__(self, client: pymongo.MongoClient[Dict[str, Any]], name: str) -> None:
        self.Client = client
        self.Database = name
        db = self.Client.get_database(name)
        user_col = db.get_collection(UserColName)
        self.UsersRepo = UsersRespo(user_col)
