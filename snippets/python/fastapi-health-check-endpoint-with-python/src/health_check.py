from fastapi import FastAPI

app = FastAPI()


async def check_database() -> dict:
    """Replace with real database ping logic."""
    # e.g. await db.execute("SELECT 1")
    return {"status": "ok"}


async def check_redis() -> dict:
    """Replace with real Redis ping logic."""
    # e.g. await redis.ping()
    return {"status": "ok"}


@app.get("/health")
async def health_check() -> dict:
    """Check all dependencies and return aggregated health status."""
    db = await check_database()
    redis = await check_redis()
    checks = {"database": db, "redis": redis}
    healthy = all(c["status"] == "ok" for c in checks.values())
    return {"status": "healthy" if healthy else "unhealthy", "checks": checks}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
