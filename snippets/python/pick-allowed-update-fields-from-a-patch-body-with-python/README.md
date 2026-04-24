# Pick Allowed Update Fields from a PATCH Body with Python

Keep only allowed top-level fields from a PATCH-style request body.

This snippet is useful when an API accepts partial updates but should only pass a fixed allowlist of fields into its update step.

- Keeps only allowed top-level keys
- Preserves explicit `null` values
- Ignores unknown fields in the returned patch

## Example

```python
from src.pick_allowed_patch_fields import pick_allowed_patch_fields

filtered_patch = pick_allowed_patch_fields(
    {
        "display_name": "Ada",
        "bio": None,
        "role": "admin",
    },
    ["display_name", "bio"],
)

print(filtered_patch)
```

Output:

```python
{"display_name": "Ada", "bio": None}
```

## Notes

- This helper only filters top-level keys.
- Missing keys stay missing, while explicit `null` values are preserved as `None`.
- If your API should reject unknown fields instead of ignoring them, keep that behavior in a separate validation step.

## Verification

Run the tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

Verified behavior covers:

- selecting allowed fields
- preserving explicit `None` values
- ignoring disallowed fields
- returning an empty result for an empty body
- rejecting a non-mapping body
- rejecting blank allowed field names

## Files

- `src/pick_allowed_patch_fields.py`
- `tests/test_pick_allowed_patch_fields.py`