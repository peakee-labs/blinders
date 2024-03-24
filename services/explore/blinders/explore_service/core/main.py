from fastapi import APIRouter
from pydantic import BaseModel

from blinders.explore_core.main import Explore, MatchInfo


def health_check_handler():
    return {"Status": "Service Healthy"}


class ExplorePostBody(BaseModel):
    name: str
    gender: str
    major: str
    native: str
    country: str
    learnings: list[str]
    interests: list[str]
    age: int


class ServiceWorker(object):
    # redis_client: redis.Redis
    core: Explore
    router: APIRouter

    # def __init__(self, redis_client: redis.Redis, explore_core: Explore) -> None:
    def __init__(self, explore_core: Explore) -> None:
        # self.redis_client = redis_client
        self.core = explore_core
        self.router = APIRouter()
        # self.init_redis_group()

    def init_route(self) -> None:
        self.router.add_api_route("/ping", health_check_handler, methods=["GET"])
        self.router.add_api_route("/explore/embed", health_check_handler, methods=["GET"])

    def embed_explore_handler(self, body: ExplorePostBody):
        matchInfor = MatchInfo(
            userId="",
            name=body.name,
            gender=body.gender,
            major=body.major,
            native=body.native,
            country=body.country,
            learnings=body.learnings,
            interests=body.interests,
            age=body.age
        )
        embed = self.core.add_user_embed(matchInfor)
        return embed

    # def init_redis_group(self):
    #     try:
    #         res = self.redis_client.xgroup_create(
    #             "match:embed", "matchGroup", "$", mkstream=True
    #         )
    #         print(res)
    #     except Exception as e:
    #         print(e)
    #         pass

    # def loop(self):
    #     consumer_name = os.getenv("REDIS_CONSUMER_NAME", "default")
    #     print("listening to stream, consumer name: ", consumer_name)
    #     while True:
    #         entries = self.redis_client.xreadgroup(
    #             "matchGroup",
    #             consumer_name,
    #             {"match:embed": ">"},
    #             block=1000,
    #             count=1,
    #         )
    #         if not isinstance(entries, list):
    #             print("reply  with unexpected format")
    #             print(entries)
    #             continue
    #
    #         if entries is None or len(entries) == 0:
    #             continue
    #
    #         user_id = entries[0][1][0][1]["id"]  # type: str
    #         if not isinstance(user_id, str):
    #             print("could not found id in stream entry")
    #             continue
    #
    #         doc = self.core.match_col.find_one({"firebaseUID": user_id})
    #         if doc is None:
    #             print("user not found")
    #             return
    #
    #         info = MatchInfo(
    #             userId=doc.get("userId"),
    #             name=doc.get("name"),
    #             gender=doc.get("gender"),
    #             major=doc.get("major"),
    #             native=doc.get("native"),
    #             country=doc.get("country"),
    #             learnings=doc.get("learnings"),
    #             interests=doc.get("interests"),
    #             age=doc.get("age"),
    #         )
    #         self.core.add_user_embed(info)
