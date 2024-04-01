import os
import time

import dotenv
import uvicorn
from fastapi import FastAPI

from blinders.embedder_core import Embedder
from blinders.embedder_service.api import API

app = FastAPI()


@app.middleware("http")
async def add_process_time_header(request, call_next):
    start_time = time.time()
    response = await call_next(request)
    process_time = time.time() - start_time
    print("{}| {}ms".format(str(time.time()), process_time))
    return response


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
    app.include_router(api.router, prefix="/embedder")
    port = os.getenv("EMBEDDER_SERVICE_PORT")
    port = int(port) if port is not None else 8084

    uvicorn.run(app, port=port)
