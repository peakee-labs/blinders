import datetime
from typing import Any, Dict

from bson.objectid import ObjectId
from pymongo.collection import Collection

from blinders.pydb.types import User, user_from_dict


class UsersRespo:
    Col: Collection[Dict[str, Any]]

    def __init__(self, col: Collection[Dict[str, Any]]) -> None:
        self.Col = col

    def insert_new_user(self, user: User) -> User:
        self.Col.insert_one(dict(user))

        return user

    def insert_raw_user(self, user: User) -> User:
        timestamp: datetime.datetime = datetime.datetime.now().replace(
            microsecond=0
        )  # since mongo have miliseconds accuracy, we could set microseconds equal to 0
        user.update(_id=ObjectId(), createdAt=timestamp, updatedAt=timestamp)

        return self.insert_new_user(user)

    def get_user_with_ID(self, id: ObjectId) -> User:
        filter = {"_id": id}
        user = self.Col.find_one(filter=filter)
        if not user:
            raise Exception("repo: user not found")
        return user_from_dict(usr=user)

    def get_user_by_firebaseUID(self, uid: str) -> User:
        filter = {"firebaseUID": uid}
        user = self.Col.find_one(filter=filter)
        if not user:
            raise Exception("repo: user not found")

        return user_from_dict(usr=user)
