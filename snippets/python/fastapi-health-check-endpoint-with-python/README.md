# FastAPI Health Check Endpoint with Python

A `/health` endpoint that checks database and Redis, then returns aggregated service status.

## What It Does

- Defines a FastAPI app with a `GET /health` endpoint
- Calls `check_database()` and `check_redis()` to probe each dependency
- Aggregates results: if all checks return `{"status": "ok"}`, the service is `healthy`; otherwise `unhealthy`
- Stub check functions are provided — replace the body with real connection logic

## Usage

```python
from health_check import app

# Run directly:
# python src/health_check.py
# Then visit http://localhost:8000/health
```

Response when all dependencies are healthy:

```json
{
  "status": "healthy",
  "checks": {
    "database": {"status": "ok"},
    "redis": {"status": "ok"}
  }
}
```

## How to Extend

Replace the stub check functions with real logic:

```python
async def check_database() -> dict:
    try:
        await db.execute("SELECT 1")
        return {"status": "ok"}
    except Exception as e:
        return {"status": "error", "detail": str(e)}
```

## Verification

```bash
cd snippets/python/fastapi-health-check-endpoint-with-python
pip install fastapi uvicorn httpx
python -m unittest discover -s tests -p "test_*.py"
```
