# Sort Tasks by Dependency Order with Python

Return a stable task order that satisfies declared task dependencies.

This snippet is useful when build steps, migration tasks, or setup jobs depend on one another and you need one safe execution order.

## Highlights

- Produces stable task order
- Detects cycles early
- Rejects unknown dependencies

## Use Cases

- Order build and packaging steps
- Run setup jobs with declared prerequisites
- Schedule migration tasks safely

## Code

```python
from src.sort_tasks_by_dependency_order import sort_tasks_by_dependency_order


result = sort_tasks_by_dependency_order(
    {
        "package": ["test"],
        "test": ["build"],
        "build": [],
    }
)
print(result)
```

## Notes

- The input maps each task name to a list of task names it depends on.
- All dependencies must also appear as keys in the input dictionary.
- Independent tasks keep their original input order when more than one valid answer exists.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- linear dependency chains
- stable ordering for independent tasks
- shared dependencies
- empty input handling
- unknown dependency errors
- duplicate dependency errors
- cyclic dependency errors

## Files

- `src/sort_tasks_by_dependency_order.py`
- `tests/test_sort_tasks_by_dependency_order.py`
- `snippet.json`