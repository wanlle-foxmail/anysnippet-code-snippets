# Load Typed Environment Settings with Python

Load a small typed settings object from environment variables with string, integer, and boolean parsing.

This snippet is useful when a script or service needs a few required and optional settings but does not need a full configuration framework.

## Highlights

- Parses text, int, and bool
- Supports sensible defaults
- Fails fast on bad values

## Use Cases

- Load app mode, port, and debug flags
- Validate required environment variables at startup
- Keep small scripts honest without a config library

## Code

```python
from src.load_typed_env_settings import load_typed_env_settings


settings = load_typed_env_settings(
    {
        "APP_ENV": "development",
        "PORT": "8000",
        "DEBUG": "false",
    }
)
print(settings)
```

## Notes

- `APP_ENV` is required.
- `PORT` defaults to `8000` when it is missing.
- `DEBUG` defaults to `False` when it is missing.
- `DEBUG` accepts only `true`, `false`, `1`, and `0`.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- loading required and optional values
- using defaults for optional settings
- rejecting missing required values
- rejecting invalid integer values
- rejecting invalid boolean values
- accepting a custom environment mapping

## Files

- `src/load_typed_env_settings.py`
- `tests/test_load_typed_env_settings.py`
- `snippet.json`