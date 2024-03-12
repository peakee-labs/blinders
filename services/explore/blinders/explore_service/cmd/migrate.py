# migrate migrate users from firestore collection to mongo db.
import datetime
import os

import dotenv
from firebase_admin import firestore, credentials, initialize_app
from pymongo import MongoClient
from bson.objectid import ObjectId
from redis.client import Redis

from blinders.explore_core.main import Explore, Embedder
from blinders.explore_service.cmd.mock import random_match_profile

userColName = "users"
matchColName = "matches"

if __name__ == "__main__":
    dotenv.load_dotenv()
    try:
        mongoURL = "mongodb://{}:{}/{}".format(
            os.getenv("MONGO_HOST"),
            os.getenv("MONGO_PORT"),
            os.getenv("MONGO_DATABASE"),
        )
        mongo_client = MongoClient(mongoURL)
        db = mongo_client[os.getenv("MONGO_DATABASE", "Default")]
        match_col = db[matchColName]
        creds = credentials.Certificate("./firebase.admin.development.json")
        app = initialize_app(creds)
        firestoreClient = firestore.client(app)
        userDocs = firestoreClient.collection("Users").stream()
        embedder = Embedder()
        redis_client = Redis(host="localhost", port=6379)
        explore = Explore(redis_client, embedder, match_col)

        for user_doc in userDocs:
            doc = user_doc.to_dict()
            if doc is None:
                userDocs.close()
                break

            firebase_UID = doc.get("firebaseUid")
            name = doc.get("name")
            image_url = doc.get("imageUrl")

            if (
                    name is None
                    or image_url is None
            ):
                continue

            now = datetime.datetime.now()
            mongo_user = db[userColName].insert_one(
                {
                    "name": name,
                    "imageURL": image_url,
                    "firebaseUID": firebase_UID,
                    "friends": [],
                    "createdAt": now,
                    "updatedAt": now,
                }
            )

            info = random_match_profile(str(mongo_user.inserted_id), name)
            match_col.insert_one({
                "userId": mongo_user.inserted_id,
                "name": info.name,
                "gender": info.gender,
                "learnings": info.learnings,
                "major": info.major,
                "native": info.native,
                "country": info.country,
                "interests": info.interests,
                "age": info.age,
            })
            explore.add_user_embed(info)

    except Exception as e:
        raise e
