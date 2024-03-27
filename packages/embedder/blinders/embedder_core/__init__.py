from sentence_transformers import SentenceTransformer


class Embedder(object):
    model: SentenceTransformer

    def __init__(self, model_name: str = "all-MiniLM-L6-v2"):
        print("loading embedder model")
        self.model = SentenceTransformer(model_name)

    def embed(self, data: str) -> list[float]:
        embeddings = self.model.encode([data])
        return [float(v) for v in embeddings[0]]
