# Refactor Large If Elif Dispatch in Python

Replace a long `if`/`elif` dispatch block with a handler registry that routes work by key.

This snippet is useful when channel names, action types, or event names decide which block of code should run and the branch list is getting hard to maintain.

## Highlights

- Replaces large branching blocks with a handler mapping
- Normalizes dispatch keys before lookup
- Supports custom handler injection without editing core dispatch logic
- Passes a payload copy into handlers to reduce accidental caller mutation

## Use Cases

- Route notifications by channel name
- Dispatch commands by action type
- Replace large `if`/`elif` blocks with a maintainable registry

## Code

```python
from src.dispatch_notification import dispatch_notification


result = dispatch_notification(
    "email",
    {
        "recipient": "dev@example.com",
        "subject": "Build finished",
        "message": "Deployment completed successfully.",
    },
)

print(result["handled_by"])
print(result["output"])
```

## Notes

- The dispatcher normalizes channel names with `strip().lower()` before lookup.
- Custom handlers can be added by passing a `handlers` mapping.
- This snippet focuses on dispatch structure, not transport-specific networking code.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- default email dispatch
- channel-specific SMS behavior
- normalized channel lookup
- unknown channel rejection
- custom handler injection
- payload-copy isolation
- blank channel validation

## Files

- `src/dispatch_notification.py`
- `tests/test_dispatch_notification.py`
- `snippet.json`