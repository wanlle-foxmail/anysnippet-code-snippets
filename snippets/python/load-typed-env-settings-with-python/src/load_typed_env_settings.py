import os
from dataclasses import dataclass
from typing import Mapping, Optional


@dataclass(frozen=True)
class AppSettings:
    app_env: str
    port: int
    debug: bool


def load_typed_env_settings(env: Optional[Mapping[str, str]] = None) -> AppSettings:
    """Load a small typed settings object from environment variables."""
    # Flow:
    #   read APP_ENV, PORT, and DEBUG from the source mapping
    #      |
    #      +-> parse text, int, and bool values -> return AppSettings
    #      `-> missing or invalid values -> raise ValueError
    source = os.environ if env is None else env

    app_env = _read_required_text(source, "APP_ENV")
    port = _read_int(source, "PORT", default=8000)
    debug = _read_bool(source, "DEBUG", default=False)

    return AppSettings(app_env=app_env, port=port, debug=debug)


def _read_required_text(source: Mapping[str, str], key: str) -> str:
    value = source.get(key)
    if not isinstance(value, str) or value.strip() == "":
        raise ValueError(f"{key} is required")
    return value.strip()


def _read_int(source: Mapping[str, str], key: str, default: int) -> int:
    value = source.get(key)
    if value is None:
        return default
    if not isinstance(value, str):
        raise ValueError(f"{key} must be a string integer")
    try:
        return int(value.strip())
    except ValueError as error:
        raise ValueError(f"{key} must be a valid integer") from error


def _read_bool(source: Mapping[str, str], key: str, default: bool) -> bool:
    value = source.get(key)
    if value is None:
        return default
    if not isinstance(value, str):
        raise ValueError(f"{key} must be a string boolean")

    normalized = value.strip().lower()
    if normalized in {"1", "true"}:
        return True
    if normalized in {"0", "false"}:
        return False

    raise ValueError(f"{key} must be one of: true, false, 1, 0")


if __name__ == "__main__":
    settings = load_typed_env_settings(
        {
            "APP_ENV": "development",
            "PORT": "8000",
            "DEBUG": "false",
        }
    )
    print(settings)