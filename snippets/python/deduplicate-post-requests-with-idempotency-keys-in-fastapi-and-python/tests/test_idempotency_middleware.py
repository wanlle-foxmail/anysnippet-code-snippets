import os
import sys
import threading
import unittest

from fastapi.responses import JSONResponse, Response
from fastapi.testclient import TestClient


SNIPPET_ROOT = os.path.join(os.path.dirname(__file__), "..", "src")
if SNIPPET_ROOT not in sys.path:
    sys.path.insert(0, SNIPPET_ROOT)

from idempotency_middleware import MemoryIdempotencyStore, new_app


class IdempotencyMiddlewareTests(unittest.TestCase):
    def test_requires_key_for_post_requests(self):
        client = TestClient(new_app())

        response = client.post("/orders")

        self.assertEqual(400, response.status_code)
        self.assertEqual("Idempotency-Key header is required", response.json()["detail"])

    def test_caches_successful_post_responses(self):
        client = TestClient(new_app())

        first_response = client.post("/orders", headers={"Idempotency-Key": "key-1"})
        second_response = client.post("/orders", headers={"Idempotency-Key": "key-1"})

        self.assertEqual(201, first_response.status_code)
        self.assertEqual(201, second_response.status_code)
        self.assertEqual("order-123", second_response.headers["x-order-id"])
        self.assertEqual(first_response.json()["order_id"], second_response.json()["order_id"])

    def test_separates_different_keys(self):
        call_count = {"value": 0}

        async def process_order() -> Response:
            call_count["value"] += 1
            order_id = "order-%d" % call_count["value"]
            return JSONResponse(
                status_code=201,
                content={"order_id": order_id},
                headers={"X-Order-ID": order_id},
            )

        client = TestClient(new_app(process_order=process_order))

        first_response = client.post("/orders", headers={"Idempotency-Key": "key-1"})
        second_response = client.post("/orders", headers={"Idempotency-Key": "key-2"})

        self.assertEqual("order-1", first_response.json()["order_id"])
        self.assertEqual("order-2", second_response.json()["order_id"])
        self.assertEqual(2, call_count["value"])

    def test_returns_conflict_while_request_is_in_progress(self):
        started = threading.Event()
        unblock = threading.Event()

        async def process_order() -> Response:
            started.set()
            await asyncio_to_thread(unblock.wait)
            return JSONResponse(status_code=201, content={"status": "created"})

        client = TestClient(new_app(store=MemoryIdempotencyStore(), process_order=process_order))
        first_response_holder = {}

        def send_first_request() -> None:
            first_response_holder["response"] = client.post("/orders", headers={"Idempotency-Key": "same-key"})

        first_request_thread = threading.Thread(target=send_first_request)
        first_request_thread.start()
        started.wait(timeout=2)

        second_response = client.post("/orders", headers={"Idempotency-Key": "same-key"})
        self.assertEqual(409, second_response.status_code)
        self.assertEqual("request already in progress", second_response.json()["detail"])

        unblock.set()
        first_request_thread.join(timeout=2)
        self.assertEqual(201, first_response_holder["response"].status_code)

    def test_does_not_cache_failed_responses(self):
        call_count = {"value": 0}

        async def process_order() -> Response:
            call_count["value"] += 1
            if call_count["value"] == 1:
                return JSONResponse(status_code=500, content={"status": "failed"})
            return JSONResponse(status_code=201, content={"status": "created"})

        client = TestClient(new_app(store=MemoryIdempotencyStore(), process_order=process_order))

        first_response = client.post("/orders", headers={"Idempotency-Key": "retry-key"})
        second_response = client.post("/orders", headers={"Idempotency-Key": "retry-key"})

        self.assertEqual(500, first_response.status_code)
        self.assertEqual(201, second_response.status_code)
        self.assertEqual("created", second_response.json()["status"])
        self.assertEqual(2, call_count["value"])

    def test_bypasses_non_post_requests(self):
        client = TestClient(new_app())

        response = client.get("/orders")

        self.assertEqual(200, response.status_code)
        self.assertEqual(b"", response.content)

    def test_separates_different_users_with_the_same_key(self):
        call_count = {"value": 0}

        async def process_order() -> Response:
            call_count["value"] += 1
            return JSONResponse(status_code=201, content={"call": call_count["value"]})

        client = TestClient(new_app(store=MemoryIdempotencyStore(), process_order=process_order))

        first_response = client.post(
            "/orders",
            headers={"Idempotency-Key": "same-key", "X-User-ID": "user-a"},
        )
        second_response = client.post(
            "/orders",
            headers={"Idempotency-Key": "same-key", "X-User-ID": "user-b"},
        )

        self.assertEqual(201, first_response.status_code)
        self.assertEqual(201, second_response.status_code)
        self.assertEqual(2, call_count["value"])

    def test_normalizes_authorization_scope(self):
        call_count = {"value": 0}

        async def process_order() -> Response:
            call_count["value"] += 1
            return JSONResponse(status_code=201, content={"call": call_count["value"]})

        client = TestClient(new_app(store=MemoryIdempotencyStore(), process_order=process_order))

        first_response = client.post(
            "/orders",
            headers={"Idempotency-Key": "same-key", "Authorization": "Bearer shared-token"},
        )
        second_response = client.post(
            "/orders",
            headers={"Idempotency-Key": "same-key", "Authorization": "bearer shared-token"},
        )

        self.assertEqual(201, first_response.status_code)
        self.assertEqual(201, second_response.status_code)
        self.assertEqual(1, call_count["value"])


async def asyncio_to_thread(function):
    import asyncio

    await asyncio.to_thread(function)


if __name__ == "__main__":
    unittest.main()