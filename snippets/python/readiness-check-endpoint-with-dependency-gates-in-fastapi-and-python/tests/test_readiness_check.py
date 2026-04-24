import os
import sys
import unittest
from unittest.mock import AsyncMock, patch

from fastapi.testclient import TestClient


SNIPPET_ROOT = os.path.join(os.path.dirname(__file__), "..", "src")
if SNIPPET_ROOT not in sys.path:
    sys.path.insert(0, SNIPPET_ROOT)

from readiness_check import app


class ReadinessCheckTests(unittest.TestCase):
    def setUp(self):
        self.client = TestClient(app)

    def test_returns_ready_when_all_required_checks_are_ready(self):
        response = self.client.get("/ready")

        self.assertEqual(200, response.status_code)
        self.assertEqual("ready", response.json()["status"])

    @patch("readiness_check.check_database_ready", new_callable=AsyncMock, return_value={"ready": False})
    def test_returns_not_ready_when_database_gate_fails(self, _mock_database):
        response = self.client.get("/ready")

        self.assertEqual(503, response.status_code)
        self.assertEqual("not_ready", response.json()["status"])
        self.assertFalse(response.json()["checks"]["database"]["ready"])

    @patch("readiness_check.check_migrations_ready", new_callable=AsyncMock, return_value={"ready": False})
    def test_returns_not_ready_when_migration_gate_fails(self, _mock_migrations):
        response = self.client.get("/ready")

        self.assertEqual(503, response.status_code)
        self.assertEqual("not_ready", response.json()["status"])
        self.assertFalse(response.json()["checks"]["migrations"]["ready"])

    @patch("readiness_check.check_cache_ready", new_callable=AsyncMock, return_value={"ready": False})
    def test_optional_cache_gate_does_not_block_readiness(self, _mock_cache):
        response = self.client.get("/ready")

        self.assertEqual(200, response.status_code)
        self.assertEqual("ready", response.json()["status"])
        self.assertFalse(response.json()["checks"]["cache"]["ready"])

    @patch("readiness_check.check_migrations_ready", new_callable=AsyncMock, return_value={"ready": False})
    @patch("readiness_check.check_database_ready", new_callable=AsyncMock, return_value={"ready": False})
    def test_multiple_required_failures_still_return_not_ready(self, _mock_database, _mock_migrations):
        response = self.client.get("/ready")

        self.assertEqual(503, response.status_code)
        self.assertFalse(response.json()["checks"]["database"]["ready"])
        self.assertFalse(response.json()["checks"]["migrations"]["ready"])

    def test_response_marks_required_and_optional_gates(self):
        response = self.client.get("/ready")
        payload = response.json()

        self.assertTrue(payload["checks"]["database"]["required"])
        self.assertTrue(payload["checks"]["migrations"]["required"])
        self.assertFalse(payload["checks"]["cache"]["required"])


if __name__ == "__main__":
    unittest.main()