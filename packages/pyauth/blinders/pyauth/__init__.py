from typing import Any, Dict, TypedDict

import firebase_admin
from firebase_admin import auth
from typing_extensions import NotRequired


class AuthUser(TypedDict):
    Email: str
    Name: str
    AuthID: str
    ID: NotRequired[str]  # hex string of users


class AuthManager(object):
    client: auth.Client
    app: firebase_admin.App

    def __init__(self, admin_json: Dict[str, Any]) -> None:
        credential = firebase_admin.credentials.Certificate(cert=admin_json)
        self.app = firebase_admin.initialize_app(credential=credential)
        self.client = auth.Client(self.app)

    def verify(self, jwt: str) -> AuthUser | None:
        try:
            authToken = self.client.verify_id_token(id_token=jwt)
        except Exception as e:
            print("auth: cannot verify id token", e)
            return None

        if authToken is None:
            return None

        firebaseUID: str = authToken.get("user_id")
        email: str = authToken.get("email")
        name: str = authToken.get("name")

        return AuthUser(Email=email, Name=name, AuthID=firebaseUID)
