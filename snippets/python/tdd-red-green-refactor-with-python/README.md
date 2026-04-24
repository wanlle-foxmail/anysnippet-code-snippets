# TDD Red Green Refactor with Python

Practice a simple RED GREEN REFACTOR loop in Python with a tiny shipping fee function and executable tests.

This snippet is useful when you want one concrete TDD example that stays small enough to follow in one sitting.

## Highlights

- Shows one RED GREEN loop
- Uses a tiny business rule
- Ends with a safe refactor

## What This Snippet Covers

- A step-by-step Markdown entry guide for the TDD loop
- One reusable Python function with visible business rules
- One unittest suite that locks in behavior and edge cases

## Entry File

- `tdd-red-green-refactor.md`
- The entry file stays in Markdown so imported users see the learning flow first.

## Example Code

```python
from src.calculate_shipping_fee import calculate_shipping_fee


print(calculate_shipping_fee(3200))
```

## Notes

- The example keeps the function intentionally small so the TDD loop stays visible.
- The tests include the Python `bool` versus `int` gotcha because `True` is an instance of `int`.
- The final refactor only lifts constants and keeps the behavior unchanged.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- default standard shipping below the threshold
- free shipping at the threshold
- premium customer free shipping
- negative subtotal rejection
- bool subtotal rejection
- unknown customer tier rejection
- boundary behavior just below the threshold

## Files

- `tdd-red-green-refactor.md`
- `src/calculate_shipping_fee.py`
- `tests/test_calculate_shipping_fee.py`
- `snippet.json`