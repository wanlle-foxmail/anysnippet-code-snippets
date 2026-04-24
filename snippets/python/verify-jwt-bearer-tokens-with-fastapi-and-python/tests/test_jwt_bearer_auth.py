import os
import sys
import unittest
from datetime import datetime, timedelta, timezone

import jwt
from fastapi.testclient import TestClient


SNIPPET_ROOT = os.path.join(os.path.dirname(__file__), "..", "src")
if SNIPPET_ROOT not in sys.path:
    sys.path.insert(0, SNIPPET_ROOT)

from jwt_bearer_auth import JWT_ALGORITHM, new_app


TEST_SECRET = "test-secret-key-for-hs256-and-hs384-validation-2026-abcdef"
WRONG_SECRET = "wrong-secret-key-for-hs256-and-hs384-validation-2026-abcdef"


def make_token(secret_key: str, claims, algorithm: str = JWT_ALGORITHM) -> str:
    return jwt.encode(claims, secret_key, algorithm=algorithm)


class JWTBearerAuthTests(unittest.TestCase):
    def setUp(self):
        self.client = TestClient(new_app(TEST_SECRET))

    def test_allows_a_valid_token(self):
        token = make_token(
            TEST_SECRET,
            {
                "sub": "user-123",
                "role": "admin",
                "exp": datetime.now(timezone.utc) + timedelta(minutes=10),
            },
        )

        response = self.client.get("/profile", headers={"Authorization": f"Bearer {token}"})

        self.assertEqual(200, response.status_code)
        self.assertEqual({"sub": "user-123", "role": "admin"}, response.json())

    def test_rejects_a_missing_authorization_header(self):
        response = self.client.get("/profile")

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_a_non_bearer_authorization_header(self):
        response = self.client.get("/profile", headers={"Authorization": "Basic abc.def.ghi"})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_a_blank_bearer_token(self):
        response = self.client.get("/profile", headers={"Authorization": "Bearer    "})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_a_malformed_jwt_string(self):
        response = self.client.get("/profile", headers={"Authorization": "Bearer not-a-jwt"})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_a_token_with_an_invalid_signature(self):
        token = make_token(
            WRONG_SECRET,
            {
                "sub": "user-123",
                "exp": datetime.now(timezone.utc) + timedelta(minutes=10),
            },
        )

        response = self.client.get("/profile", headers={"Authorization": f"Bearer {token}"})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_an_expired_token(self):
        token = make_token(
            TEST_SECRET,
            {
                "sub": "user-123",
                "exp": datetime.now(timezone.utc) - timedelta(minutes=10),
            },
        )

        response = self.client.get("/profile", headers={"Authorization": f"Bearer {token}"})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_a_token_without_an_exp_claim(self):
        token = make_token(TEST_SECRET, {"sub": "user-123"})

        response = self.client.get("/profile", headers={"Authorization": f"Bearer {token}"})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])

    def test_rejects_an_unexpected_signing_method(self):
        token = make_token(
            TEST_SECRET,
            {
                "sub": "user-123",
                "exp": datetime.now(timezone.utc) + timedelta(minutes=10),
            },
            algorithm="HS384",
        )

        response = self.client.get("/profile", headers={"Authorization": f"Bearer {token}"})

        self.assertEqual(401, response.status_code)
        self.assertEqual("invalid bearer token", response.json()["detail"])


if __name__ == "__main__":
    unittest.main()