import json
from typing import Any, Dict

from blinders.pyauth import AuthManager, AuthUser
from blinders.pyauth.aws import auth
from blinders.pydb import MongoManager
from blinders.pydb.types import User
from blinders.pydb.utils.mongo import init_mongo_client

testURL = "mongodb://localhost:27017"
testDB = "blinders-test"


def test_aws():
    auth_cert: Dict[str, Any] = {}
    with open("firebase.admin.json") as file:
        auth_cert = json.load(file)

    auth_manager = AuthManager(auth_cert)
    valid_firebaseuid = "t7ZYtyjYCbMxOefUALu8b2P4AVO2"
    user = User(
        name="name", firebaseUID=valid_firebaseuid, email="email", friends=[], imageURL="imageURL"
    )
    mongo_manager = MongoManager(
        client=init_mongo_client(testURL),
        name=testDB,
    )

    added_user = mongo_manager.UsersRepo.insert_raw_user(user=user)
    mockEvent = {"headers": {"authorization": "Bearer sampletoken"}}

    def test_handler(auth_user: AuthUser, event=mockEvent):
        assert auth_user.get("AuthID") == valid_firebaseuid
        assert auth_user.get("ID") == str(added_user.get("_id"))

    auth(auth_manager=auth_manager, repo=mongo_manager.UsersRepo)(test_handler)(event=mockEvent)
    mongo_manager.UsersRepo.Col.drop()
