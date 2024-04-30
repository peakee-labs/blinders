from blinders.pyauth import AuthManager
from blinders.pydb.repository.users import UsersRespo
from typing import Dict, Any


default_headers = {
    "Access-Control-Allow-Origin": "*",
    "Content-Type": "application/json",
}


def auth(auth_manager: AuthManager, repo: UsersRespo):
    """
    middleware try to parse and verify the incomming apigateway event header field and get user
    information

    authenticate's information then will pass to the lambda handler with param auth_user

    In order to use auth middleware, wrap this decorater to the lambda handler, and get the authUser
    inside the lambda via auth_user param.

    [DEPRECATED] This function is deprecated due to its dependencies significantly increasing
    the size of the bundle.

    Lambda need to authenticate the request with bearer token now should invoke authenticate
    function with pytransport package
    """

    def aws_auth(handler):
        # this method try to authorized event from api gateway using
        # event.header.Authorization will be used
        # event from inter system service will not need to authorized
        def wrapper(**kwargs: Dict[str, Any]):
            event: Dict[str, Any] = kwargs["event"]
            if event is None:
                raise Exception("event field not exsisted in lambda params")
                # try to authorized with this header
            headers: Dict[str, Any] = event["headers"]
            if headers is None:
                raise Exception("headers field not exsisted in event")

            auth: str = headers["authorization"]
            if auth is None:
                return {
                    "statusCode": 400,
                    "body": "missing authorization header",
                    "headers": default_headers,
                }
            if not auth.startswith("Bearer "):
                return {
                    "statusCode": 400,
                    "body": "invalid jwt, missing bearer token",
                    "headers": default_headers,
                }
            token = auth.split(" ")[1]
            userAuth = auth_manager.verify(token)
            if userAuth is None:
                return {
                    "statusCode": 400,
                    "body": "invalid jwt, token cannot verify given token",
                    "headers": default_headers,
                }

            userRepo = repo.get_user_by_firebaseUID(userAuth.get("AuthID"))
            userID = userRepo.get("_id")
            if userID:
                userAuth["ID"] = str(userID)

            return handler(auth_user=userAuth, **kwargs)

        return wrapper

    return aws_auth
