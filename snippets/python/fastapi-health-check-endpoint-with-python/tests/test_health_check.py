import os
import sys
import unittest
from unittest.mock import AsyncMock, patch

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "src"))

from fastapi.testclient import TestClient

from health_check import app


class TestHealthCheck(unittest.TestCase):
    def setUp(self):
        self.client = TestClient(app)

    def test_all_healthy(self):
        resp = self.client.get("/health")
        self.assertEqual(resp.status_code, 200)
        data = resp.json()
        self.assertEqual(data["status"], "healthy")
        self.assertEqual(data["checks"]["database"]["status"], "ok")
        self.assertEqual(data["checks"]["redis"]["status"], "ok")

    @patch("health_check.check_database", new_callable=AsyncMock, return_value={"status": "error"})
    def test_database_failure(self, _mock_db):
        resp = self.client.get("/health")
        data = resp.json()
        self.assertEqual(data["status"], "unhealthy")
        self.assertEqual(data["checks"]["database"]["status"], "error")
        self.assertEqual(data["checks"]["redis"]["status"], "ok")

    @patch("health_check.check_redis", new_callable=AsyncMock, return_value={"status": "error"})
    def test_redis_failure(self, _mock_redis):
        resp = self.client.get("/health")
        data = resp.json()
        self.assertEqual(data["status"], "unhealthy")
        self.assertEqual(data["checks"]["redis"]["status"], "error")
        self.assertEqual(data["checks"]["database"]["status"], "ok")

    @patch("health_check.check_redis", new_callable=AsyncMock, return_value={"status": "error"})
    @patch("health_check.check_database", new_callable=AsyncMock, return_value={"status": "error"})
    def test_all_dependencies_down(self, _mock_db, _mock_redis):
        resp = self.client.get("/health")
        data = resp.json()
        self.assertEqual(data["status"], "unhealthy")
        self.assertEqual(data["checks"]["database"]["status"], "error")
        self.assertEqual(data["checks"]["redis"]["status"], "error")

    def test_response_status_code_is_200(self):
        resp = self.client.get("/health")
        self.assertEqual(resp.status_code, 200)

    def test_response_structure(self):
        resp = self.client.get("/health")
        data = resp.json()
        self.assertIn("status", data)
        self.assertIn("checks", data)
        self.assertIn("database", data["checks"])
        self.assertIn("redis", data["checks"])


if __name__ == "__main__":
    unittest.main()
