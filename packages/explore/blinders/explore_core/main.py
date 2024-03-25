from blinders.explore_core.embedder import Embedder
from blinders.explore_core.types import MatchInfo


class Explore(object):
    embedder: Embedder

    def __init__(self, embedder: Embedder) -> None:
        self.embedder = embedder

    def add_user_embed(self, info: MatchInfo) -> list[float]:
        """
        add_use_embed call after a new match entry already added to matches collection, this will embed recently
        document then add to vector database.
        :param info: blinders.explore_core.types.MatchInfo
        """
        return self.embedder.embed(info)
