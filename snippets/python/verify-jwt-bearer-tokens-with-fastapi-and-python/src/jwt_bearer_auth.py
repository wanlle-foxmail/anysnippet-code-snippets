from typing import Any, Callable, Dict

import jwt
from fastapi import Depends, FastAPI, Header, HTTPException


JWT_ALGORITHM = "HS256"
DEMO_SECRET_KEY = "demo-secret-key-for-hs256-example-2026"


def extract_bearer_token(authorization_header: str) -> str:
    """Return the bearer token from an Authorization header."""
    parts = authorization_header.strip().split()
    if len(parts) != 2 or parts[0].lower() != "bearer":
        raise ValueError("authorization header must use bearer scheme")
    if parts[1].strip() == "":
        raise ValueError("bearer token is required")
    return parts[1]


def verify_jwt_bearer_token(authorization_header: str, secret_key: str) -> Dict[str, Any]:
    """Validate one HS256 bearer token and return its claims."""
    # Flow:
    #   read Authorization header
    #      |
    #      +-> missing or malformed bearer token -> raise ValueError
    #      `-> decode HS256 token -> require exp -> return claims
    token_string = extract_bearer_token(authorization_header)
    try:
        claims = jwt.decode(
            token_string,
            secret_key,
            algorithms=[JWT_ALGORITHM],
            options={"require": ["exp"]},
        )
    except jwt.InvalidTokenError as error:
        raise ValueError("invalid bearer token") from error

    return claims


def build_jwt_dependency(secret_key: str) -> Callable[..., Dict[str, Any]]:
    """Build one FastAPI dependency that returns validated JWT claims."""

    def require_jwt_claims(authorization: str = Header(default="", alias="Authorization")) -> Dict[str, Any]:
        try:
            return verify_jwt_bearer_token(authorization, secret_key)
        except ValueError as error:
            raise HTTPException(status_code=401, detail="invalid bearer token") from error

    return require_jwt_claims


def new_app(secret_key: str) -> FastAPI:
    """Create one FastAPI app with a protected profile route."""
    app = FastAPI()
    require_jwt_claims = build_jwt_dependency(secret_key)

    @app.get("/profile")
    async def profile(claims: Dict[str, Any] = Depends(require_jwt_claims)) -> Dict[str, str]:
        subject = claims.get("sub")
        role = claims.get("role")
        return {
            "sub": subject if isinstance(subject, str) else "",
            "role": role if isinstance(role, str) else "",
        }

    return app


def make_example_token(secret_key: str) -> str:
    """Create one short-lived example bearer token."""
    return jwt.encode(
        {
            "sub": "demo-user",
            "role": "admin",
            "exp": 4102445700,
        },
        secret_key,
        algorithm=JWT_ALGORITHM,
    )


app = new_app(DEMO_SECRET_KEY)


if __name__ == "__main__":
    import uvicorn

    print(make_example_token(DEMO_SECRET_KEY))
    uvicorn.run(app, host="0.0.0.0", port=8000)