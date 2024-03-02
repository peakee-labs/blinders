from typing_extensions import IntVar
import redis

import os
from blinders.explore_core.main import Explore
from blinders.explore_core.types import MatchInfo


class ServiceWorker(object):
    redisClient: redis.Redis
    core: Explore

    def __init__(self, redisClient: redis.Redis, exploreCore: Explore) -> None:
        self.redisClient = redisClient
        self.core = exploreCore
        self.initRedisGroup()

    def initRedisGroup(self):
        try:
            res = self.redisClient.xgroup_create(
                "match:embed", "matchGroup", "$", mkstream=True
            )
            print(res)
        except Exception as e:
            print(e)
            pass

    def loop(self):
        print("looping")
        i = 0
        while True:
            entries = self.redisClient.xreadgroup(
                "matchGroup",
                os.getenv("REDIS_CONSUMER", "default"),
                {"match:embed": ">"},
                block=10000,
            )
            if not isinstance(entries, list):
                print("reply  with unexpected format")
                print(entries)
                continue

            if entries is None or len(entries) == 0:
                print("timeout tried again")
                continue

            userID = entries[0][1][0][1]["id"]  # type: str
            if not isinstance(userID, str):
                print("could not found id in stream entry")
                continue

            doc = self.core.matchCol.find_one({"firebaseUID": userID})
            if doc is None:
                print("user not found")
                return

            info = MatchInfo(
                firebaseUID=userID,
                name=doc.get("name"),
                gender=doc.get("gender"),
                major=doc.get("major"),
                native=doc.get("native"),
                country=doc.get("country"),
                learnings=doc.get("learnings"),
                interests=doc.get("interests"),
                userID=doc.get("userID"),
                age=doc.get("age"),
            )
            self.core.AddUserEmbed(info)
            i += 1
