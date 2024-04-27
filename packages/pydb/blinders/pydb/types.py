from typing import TypedDict, List, Any, Dict

# with python version > 3.11, NotRequired is in typing package
from typing_extensions import NotRequired

from bson import ObjectId
import datetime


class User(TypedDict):
    _id: NotRequired[ObjectId]
    name: str
    email: str
    firebaseUID: str
    imageURL: str
    friends: List[ObjectId]
    createdAt: NotRequired[datetime.datetime]
    updatedAt: NotRequired[datetime.datetime]


def user_from_dict(usr: Dict[str, Any]) -> User:
    return User(
        _id=ObjectId(usr.get("_id", "")),
        name=usr.get("name", ""),
        email=usr.get("email", ""),
        firebaseUID=usr.get("firebaseUID", ""),
        imageURL=usr.get("imageURL", ""),
        friends=usr.get("friends", []),
        createdAt=usr.get("createdAt", datetime.datetime.now()),
        updatedAt=usr.get("updatedAt", datetime.datetime.now()),
    )
