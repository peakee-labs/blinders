# from bson.objectid import ObjectId
# from pymongo.collection import Collection
# from redis.client import Redis
# from redis.commands.search.field import (
#     TextField,
#     VectorField,
# )
# from redis.commands.search.indexDefinition import IndexDefinition, IndexType

from blinders.explore_core.embedder import Embedder
from blinders.explore_core.types import MatchInfo
# from blinders.explore_core.utils import create_redis_match_key


class Explore(object):
    # redis_client: Redis
    embedder: Embedder
    # match_col: Collection
    # vector_dimension = 384

    def __init__(
            self,
            # redis_client: Redis,
            embedder: Embedder,
            # match_col: Collection,
    ) -> None:
        # if not redis_client.ping():
        #     raise Exception("cannot ping to redis")

        # self.redis_client = redis_client
        self.embedder = embedder
        # self.match_col = match_col
        # self.init_redis_index()

    # def init_redis_index(self):
    #     try:
    #         schema = (
    #             TextField("$.id", no_stem=True, as_name="id"),
    #             VectorField(
    #                 "$.embed",
    #                 "HNSW",
    #                 {
    #                     "TYPE": "FLOAT32",
    #                     "DIM": self.vector_dimension,
    #                     "DISTANCE_METRIC": "L2",
    #                 },
    #                 as_name="embed",
    #             ),
    #         )
    #         definition = IndexDefinition(prefix=["match:"], index_type=IndexType.JSON)
    #         res = self.redis_client.ft("idx:match_vss").create_index(
    #             fields=schema, definition=definition
    #         )
    #         print(res)
    #
    #     except Exception as e:
    #         # maybe index defined
    #         print("error=>", e)

    def add_user_embed(self, info: MatchInfo) -> list[float]:
        """
        add_use_embed call after a new match entry already added to matches collection, this will embed recently
        document then add to vector database.
        :param info: blinders.explore_core.types.MatchInfo
        """
        # doc = self.match_col.find({"userId": ObjectId(info.userId)})
        # if doc is None:
        #     raise Exception("user not existed")
        #
        return self.embedder.embed(info)
        # self.redis_client.json().set(
        #     create_redis_match_key(info.userId),
        #     "$",
        #     {
        #         "embed": embed,
        #         "id": info.userId,
        #     },
        # )
