import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.load_typed_env_settings import AppSettings, load_typed_env_settings


class LoadTypedEnvSettingsTests(unittest.TestCase):
    def test_loads_required_and_optional_values(self):
        settings = load_typed_env_settings(
            {
                "APP_ENV": "production",
                "PORT": "8080",
                "DEBUG": "true",
            }
        )

        self.assertEqual(AppSettings(app_env="production", port=8080, debug=True), settings)

    def test_uses_defaults_for_missing_optional_values(self):
        settings = load_typed_env_settings({"APP_ENV": "development"})

        self.assertEqual(AppSettings(app_env="development", port=8000, debug=False), settings)

    def test_raises_value_error_when_required_value_is_missing(self):
        with self.assertRaises(ValueError) as raised:
            load_typed_env_settings({})

        self.assertIn("APP_ENV", str(raised.exception))

    def test_raises_value_error_for_invalid_port(self):
        with self.assertRaises(ValueError) as raised:
            load_typed_env_settings({"APP_ENV": "test", "PORT": "abc"})

        self.assertIn("PORT", str(raised.exception))

    def test_raises_value_error_for_invalid_debug_value(self):
        with self.assertRaises(ValueError) as raised:
            load_typed_env_settings({"APP_ENV": "test", "DEBUG": "sometimes"})

        self.assertIn("DEBUG", str(raised.exception))

    def test_accepts_custom_mapping_without_reading_os_environ(self):
        settings = load_typed_env_settings(
            {
                "APP_ENV": "staging",
                "PORT": "9000",
                "DEBUG": "0",
            }
        )

        self.assertEqual(AppSettings(app_env="staging", port=9000, debug=False), settings)


if __name__ == "__main__":
    unittest.main()