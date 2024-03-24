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
        self.router.add_api_route("/explore/embed", self.embed_explore_handler, methods=["POST"])

    def embed_explore_handler(self, body: ExplorePostBody):
        print(body)
        match_info = MatchInfo(
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
        embed = self.core.add_user_embed(match_info)
        return {"embed": embed}