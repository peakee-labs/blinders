import os

import dotenv
import uvicorn
from fastapi import FastAPI

from blinders.embedder_core import Embedder
from blinders.embedder_service.api import API


def main():
    env = os.getenv("ENVIRONMENT")
    env_file = ".env"
    if env != "":
        env_file = ".env.{}".format(env)

    print(
        "running embedder service on {}, envfile: {}".format(
            "production" if env is None else env, env_file
        )
    )
    dotenv.load_dotenv(env_file)

    embedder = Embedder()
    api = API(embedder)
    api.init_route()

    app = FastAPI()
    app.include_router(api.router, prefix="/embedder")
    port = os.getenv("EMBEDDER_SERVICE_PORT")
    port = int(port) if port is not None else 8084

    uvicorn.run(app, port=port)
