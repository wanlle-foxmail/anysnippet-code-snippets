# Refactor Nested If Else to a Config Map in Python

Replace nested `if` and `else` decision logic with a config mapping for support ticket routing.

This snippet is useful when multiple input dimensions decide a result and the original code is becoming a maze of nested `if` and `else` blocks.

## Highlights

- Moves rules into config
- Uses tuple keys for lookup
- Falls back to a default route

## Use Cases

- Route support tickets by channel, customer tier, and issue type
- Replace deeply nested decision trees with a reviewable rule table
- Add new business rules by editing data instead of control flow

## Code

```python
from src.route_support_ticket import route_support_ticket


result = route_support_ticket("email", "pro", "login")
print(result["queue"])
print(result["sla_hours"])
```

## Notes

- The tuple key matches `(channel, customer_tier, issue_type)`.
- `DEFAULT_ROUTE` is used when no exact rule matches.
- `dict(route)` returns a copy so callers do not mutate shared config.
- This pattern works best when rules are lookup-oriented rather than algorithm-oriented.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- enterprise billing routing
- pro login routing
- enterprise security phone routing
- default fallback selection
- rule-copy isolation

## Files

- `src/route_support_ticket.py`
- `tests/test_route_support_ticket.py`
- `snippet.json`