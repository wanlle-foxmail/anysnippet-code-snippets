import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.sort_tasks_by_dependency_order import sort_tasks_by_dependency_order


class SortTasksByDependencyOrderTests(unittest.TestCase):
    def test_orders_linear_dependencies(self):
        tasks = {
            "package": ["test"],
            "test": ["build"],
            "build": [],
        }

        ordered_tasks = sort_tasks_by_dependency_order(tasks)

        self.assertEqual(["build", "test", "package"], ordered_tasks)

    def test_preserves_input_order_for_independent_tasks(self):
        tasks = {
            "lint": [],
            "build": [],
            "docs": [],
        }

        ordered_tasks = sort_tasks_by_dependency_order(tasks)

        self.assertEqual(["lint", "build", "docs"], ordered_tasks)

    def test_handles_shared_dependencies(self):
        tasks = {
            "deploy": ["test"],
            "test": ["build"],
            "build": [],
            "lint": [],
        }

        ordered_tasks = sort_tasks_by_dependency_order(tasks)

        self.assertEqual(["build", "lint", "test", "deploy"], ordered_tasks)

    def test_returns_empty_list_for_empty_input(self):
        self.assertEqual([], sort_tasks_by_dependency_order({}))

    def test_raises_value_error_for_unknown_dependency(self):
        with self.assertRaises(ValueError) as raised:
            sort_tasks_by_dependency_order({"test": ["build"]})

        self.assertIn("unknown dependency", str(raised.exception))

    def test_raises_value_error_for_duplicate_dependency(self):
        with self.assertRaises(ValueError) as raised:
            sort_tasks_by_dependency_order({"test": ["build", "build"], "build": []})

        self.assertIn("duplicate dependency", str(raised.exception))

    def test_raises_value_error_for_cyclic_dependency(self):
        with self.assertRaises(ValueError) as raised:
            sort_tasks_by_dependency_order({"build": ["test"], "test": ["build"]})

        self.assertIn("cyclic dependency", str(raised.exception))


if __name__ == "__main__":
    unittest.main()