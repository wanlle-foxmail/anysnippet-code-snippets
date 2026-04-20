# Refactor If Elif Dispatch to a Handler Map in Python

Replace a long `if`/`elif` dispatch chain with a handler mapping for email, SMS, and webhook notifications.

This snippet is useful when channel names, action types, or event names decide which block of code should run and the branch list is getting hard to maintain.

## Highlights

- Maps channels to handlers
- Keeps handler logic separate
- Uses plain dict lookup

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

- Add new channels by putting another function in the `handlers` mapping.
- Each handler receives the payload for one channel.
- This snippet focuses on dispatch structure, not transport-specific networking code.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- email dispatch
- default email subject handling
- SMS dispatch
- webhook dispatch
- unknown channel rejection

## Files

- `src/dispatch_notification.py`
- `tests/test_dispatch_notification.py`
- `snippet.json`