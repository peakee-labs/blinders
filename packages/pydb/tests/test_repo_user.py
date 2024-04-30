from typing import Dict, Any

from pymongo import MongoClient

import blinders.pydb.types as pydb_types
import blinders.pydb as pydb

mongo_test_url = "mongodb://localhost:27017/blinders"
mongo_test_db = "blinders-test"
manager: pydb.MongoManager | None = None


def test_insert_new_user():
    manager = get_manager()
    user_repo = manager.UsersRepo
    # make sure col is empty before run test
    user_repo.Col.drop()

    firebase_UID = "firebaseUID"
    user = pydb_types.User(
        name="name",
        firebaseUID=firebase_UID,
        email="email",
        friends=[],
        imageURL="imageURL",
    )
    added_user = user_repo.insert_raw_user(user)

    user_id = added_user.get("_id")
    assert user_id is not None and user_id.is_valid
    assert added_user.get("createdAt") is not None
    assert added_user.get("updatedAt") is not None

    new_added_user = user_repo.get_user_with_ID(user_id)
    assert new_added_user == added_user

    find_with_firebaseuid = user_repo.get_user_by_firebaseUID(firebase_UID)
    assert find_with_firebaseuid == added_user


def get_manager() -> pydb.MongoManager:
    global manager
    if manager is None:
        client: MongoClient[Dict[str, Any]] = MongoClient(host=mongo_test_url)
        manager = pydb.MongoManager(client=client, name=mongo_test_db)

    return manager
