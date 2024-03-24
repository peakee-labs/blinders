import os

import dotenv
import uvicorn
from fastapi import FastAPI

from blinders.explore_core.embedder import Embedder
from blinders.explore_core.main import Explore
from blinders.explore_service.core.main import ServiceWorker

if __name__ == "__main__":
    try:
        env = os.getenv("ENVIRONMENT")
        env_file = ".env"
        if env != "":
            env_file = ".env.{}".format(env)
        print("running embedder service on {}, envfile: {}".format("production" if env is None else env, env_file))
        dotenv.load_dotenv(env_file)

        embedder = Embedder()
        explore = Explore(embedder)
        service_core = ServiceWorker(explore_core=explore)
        service_core.init_route()

        app = FastAPI()
        app.include_router(service_core.router, prefix="/api")
        port = os.getenv("EXPLORE_EMBEDDER_PORT")
        port = int(port) if port is not None else 8084

        uvicorn.run(app, host="0.0.0.0", port=port)

    except Exception as e:
        print("exception: ", e)
