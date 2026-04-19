# Replace Nested If Else with Config Mapping in Python

Replace nested branching rules with a configuration mapping that is easier to extend and review.

This snippet is useful when multiple input dimensions decide a result and the original code is becoming a maze of nested `if` and `else` blocks.

## Highlights

- Moves routing rules into a data mapping instead of nested branching
- Normalizes inputs before matching
- Supports exact matches, wildcard matches, and a default fallback
- Returns a copy of the selected rule so callers cannot mutate shared config

## Use Cases

- Route support tickets by channel, customer tier, and issue type
- Replace deeply nested decision trees with a reviewable rule table
- Add new business rules by editing data instead of control flow

## Code

```python
from src.route_support_ticket import route_support_ticket


result = route_support_ticket("chat", "enterprise", "billing")
print(result["queue"])
print(result["sla_hours"])
```

## Notes

- Inputs are normalized with `strip().lower()` before rule lookup.
- Rule precedence is `exact match` first, then wildcard combinations, then the default fallback.
- This pattern works best when rules are lookup-oriented rather than algorithm-oriented.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- exact rule selection
- issue wildcard selection
- cross-channel security routing
- default fallback selection
- input normalization
- blank input validation
- rule-copy isolation

## Files

- `src/route_support_ticket.py`
- `tests/test_route_support_ticket.py`
- `snippet.json`