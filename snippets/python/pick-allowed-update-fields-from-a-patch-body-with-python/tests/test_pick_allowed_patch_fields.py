import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.pick_allowed_patch_fields import pick_allowed_patch_fields


class PickAllowedPatchFieldsTests(unittest.TestCase):
    def test_keeps_allowed_fields(self):
        filtered_patch = pick_allowed_patch_fields(
            {
                "display_name": "Ada",
                "timezone": "UTC",
                "role": "admin",
            },
            ["display_name", "timezone"],
        )

        self.assertEqual(
            {
                "display_name": "Ada",
                "timezone": "UTC",
            },
            filtered_patch,
        )

    def test_preserves_explicit_null_values(self):
        filtered_patch = pick_allowed_patch_fields(
            {
                "bio": None,
                "role": "admin",
            },
            ["bio"],
        )

        self.assertEqual({"bio": None}, filtered_patch)

    def test_ignores_disallowed_fields(self):
        filtered_patch = pick_allowed_patch_fields(
            {
                "role": "admin",
                "is_staff": True,
            },
            ["display_name", "bio"],
        )

        self.assertEqual({}, filtered_patch)

    def test_returns_empty_mapping_for_empty_patch_body(self):
        filtered_patch = pick_allowed_patch_fields({}, ["display_name"])

        self.assertEqual({}, filtered_patch)

    def test_raises_type_error_for_non_mapping_patch_body(self):
        with self.assertRaises(TypeError) as raised:
            pick_allowed_patch_fields([("display_name", "Ada")], ["display_name"])

        self.assertIn("patch_body", str(raised.exception))

    def test_raises_value_error_for_blank_allowed_field_name(self):
        with self.assertRaises(ValueError) as raised:
            pick_allowed_patch_fields({"display_name": "Ada"}, ["display_name", " "])

        self.assertIn("allowed field names", str(raised.exception))


if __name__ == "__main__":
    unittest.main()