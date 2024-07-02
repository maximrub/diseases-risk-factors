import requests

from thesis_diseases_risk_factors.infrastructure.resources.resource import Resource
from thesis_diseases_risk_factors.graph.client import Client
import httpx


class GraphClientResource(Resource):
    @classmethod
    def _create(cls, config) -> Client:
        payload = {
            "grant_type": config["auth0"]["grantType"],
            "client_id": config["auth0"]["clientID"],
            "client_secret": config["auth0"]["clientSecret"],
            "audience": config["server"]["url"]
        }

        headers = {
            "content-type": "application/json"
        }

        response = requests.post(config["auth0"]["url"], json=payload, headers=headers)
        access_token = response.json()["access_token"]
        token_type = response.json()["token_type"]

        client = Client(config["server"]["url"], {"Authorization": f"{token_type} {access_token}"}, 
                        http_client=httpx.AsyncClient(timeout=httpx.Timeout(30.0), headers={"Authorization": f"{token_type} {access_token}"}))
        return client

    @classmethod
    def _shutdown(cls, resource: Client) -> None:
        pass