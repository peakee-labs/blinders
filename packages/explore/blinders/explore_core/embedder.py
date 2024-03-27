from sentence_transformers import SentenceTransformer

from blinders.explore_core.types import MatchInfo


class Embedder(object):
    model: SentenceTransformer

    def __init__(self, model_name: str = "all-MiniLM-L6-v2"):
        print("loading embedder model")
        self.model = SentenceTransformer(model_name)

    def embed(self, info: MatchInfo) -> list[float]:
        embed_string = str(info)
        embeddings = self.model.encode([embed_string])
        return [float(v) for v in embeddings[0]]
