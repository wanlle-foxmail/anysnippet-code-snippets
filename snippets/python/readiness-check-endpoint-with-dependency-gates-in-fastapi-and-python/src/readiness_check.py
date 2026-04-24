from typing import Dict

from fastapi import FastAPI
from fastapi.responses import JSONResponse


app = FastAPI()


async def check_database_ready() -> Dict[str, object]:
    """Replace with real database readiness logic."""
    return {"ready": True}


async def check_migrations_ready() -> Dict[str, object]:
    """Replace with real migration readiness logic."""
    return {"ready": True}


async def check_cache_ready() -> Dict[str, object]:
    """Replace with real cache readiness logic."""
    return {"ready": True}


@app.get("/ready")
async def readiness_check() -> JSONResponse:
    """Gate readiness on required checks and report optional checks separately."""
    # Flow:
    #   run required and optional readiness checks
    #      |
    #      +-> all required checks ready -> return 200 ready
    #      `-> any required check not ready -> return 503 not_ready
    checks = {
        "database": {"required": True, **(await check_database_ready())},
        "migrations": {"required": True, **(await check_migrations_ready())},
        "cache": {"required": False, **(await check_cache_ready())},
    }

    is_ready = all(check["ready"] for check in checks.values() if check["required"])
    status_code = 200 if is_ready else 503
    payload = {
        "status": "ready" if is_ready else "not_ready",
        "checks": checks,
    }
    return JSONResponse(status_code=status_code, content=payload)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)