import hashlib
from dataclasses import dataclass
from threading import Lock
from typing import Awaitable, Callable, Dict, Mapping, Optional

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse, Response


@dataclass(frozen=True)
class CachedResponse:
    status_code: int
    headers: Dict[str, str]
    body: bytes


@dataclass(frozen=True)
class _MemoryEntry:
    response: Optional[CachedResponse]
    in_progress: bool
    completed: bool


class MemoryIdempotencyStore:
    def __init__(self) -> None:
        self._lock = Lock()
        self._entries: Dict[str, _MemoryEntry] = {}

    def load(self, key: str) -> Optional[CachedResponse]:
        with self._lock:
            entry = self._entries.get(key)
            if entry is None or not entry.completed or entry.response is None:
                return None
            return CachedResponse(
                status_code=entry.response.status_code,
                headers=dict(entry.response.headers),
                body=bytes(entry.response.body),
            )

    def try_start(self, key: str) -> bool:
        with self._lock:
            entry = self._entries.get(key)
            if entry is not None and (entry.in_progress or entry.completed):
                return False
            self._entries[key] = _MemoryEntry(response=None, in_progress=True, completed=False)
            return True

    def save(self, key: str, response: CachedResponse) -> None:
        with self._lock:
            self._entries[key] = _MemoryEntry(
                response=CachedResponse(
                    status_code=response.status_code,
                    headers=dict(response.headers),
                    body=bytes(response.body),
                ),
                in_progress=False,
                completed=True,
            )

    def delete(self, key: str) -> None:
        with self._lock:
            self._entries.pop(key, None)


ProcessOrder = Callable[[], Awaitable[Response]]


def new_app(
    store: Optional[MemoryIdempotencyStore] = None,
    process_order: Optional[ProcessOrder] = None,
) -> FastAPI:
    app = FastAPI()
    active_store = MemoryIdempotencyStore() if store is None else store
    process_order_handler = default_process_order if process_order is None else process_order

    @app.middleware("http")
    async def idempotency_middleware(request: Request, call_next: Callable[[Request], Awaitable[Response]]) -> Response:
        if request.method != "POST":
            return await call_next(request)

        idempotency_key = request.headers.get("Idempotency-Key", "").strip()
        if idempotency_key == "":
            return JSONResponse(status_code=400, content={"detail": "Idempotency-Key header is required"})

        scoped_key = scoped_idempotency_key(request, idempotency_key)
        cached_response = active_store.load(scoped_key)
        if cached_response is not None:
            return response_from_cached(cached_response)
        if not active_store.try_start(scoped_key):
            cached_response = active_store.load(scoped_key)
            if cached_response is not None:
                return response_from_cached(cached_response)
            return JSONResponse(status_code=409, content={"detail": "request already in progress"})

        response = await call_next(request)
        captured_response = await capture_response(response)
        if 200 <= captured_response.status_code < 300:
            active_store.save(scoped_key, captured_response)
        else:
            active_store.delete(scoped_key)

        return response_from_cached(captured_response)

    @app.post("/orders")
    async def create_order() -> Response:
        return await process_order_handler()

    @app.get("/orders")
    async def list_orders() -> Response:
        return Response(status_code=200)

    return app


async def default_process_order() -> Response:
    return JSONResponse(
        status_code=201,
        content={"order_id": "order-123", "status": "created"},
        headers={"X-Order-ID": "order-123"},
    )


def scoped_idempotency_key(request: Request, idempotency_key: str) -> str:
    return ":".join(
        [
            request.method,
            request.url.path or "/",
            request_caller_scope(request.headers),
            idempotency_key,
        ]
    )


def request_caller_scope(headers: Mapping[str, str]) -> str:
    user_id = headers.get("X-User-ID", "").strip()
    if user_id != "":
        return "user:" + user_id

    authorization_header = headers.get("Authorization", "").strip()
    if authorization_header != "":
        return "auth:" + authorization_scope_hash(authorization_header)

    return "anonymous"


def authorization_scope_hash(authorization_header: str) -> str:
    parts = authorization_header.split()
    normalized_value = authorization_header
    if len(parts) > 0:
        normalized_value = parts[0].lower()
        if len(parts) > 1:
            normalized_value += " " + " ".join(parts[1:])
    return hashlib.sha256(normalized_value.encode("utf-8")).hexdigest()[:32]


async def capture_response(response: Response) -> CachedResponse:
    body = b""
    async for chunk in response.body_iterator:
        body += chunk

    return CachedResponse(
        status_code=response.status_code,
        headers={key: value for key, value in response.headers.items()},
        body=body,
    )


def response_from_cached(cached_response: CachedResponse) -> Response:
    headers = {key: value for key, value in cached_response.headers.items() if key.lower() != "content-length"}
    return Response(content=cached_response.body, status_code=cached_response.status_code, headers=headers)


app = new_app()


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)