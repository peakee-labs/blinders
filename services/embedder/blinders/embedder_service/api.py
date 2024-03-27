from fastapi import APIRouter

from blinders.embedder_core import Embedder


class API(object):
    embedder: Embedder
    router: APIRouter

    def __init__(self, embedder: Embedder) -> None:
        self.router = APIRouter()
        self.embedder = embedder

    def init_route(self) -> None:
        self.router.add_api_route("/ping", self.ping_handler, methods=["GET"])
        self.router.add_api_route("/embed", self.embed_handler, methods=["POST"])

    def ping_handler(self):
        return {"message": "pong"}

    def embed_handler(self, body: dict[str, str]):
        embedded = self.embedder.embed(body["data"])
        return {"embedded": embedded}
