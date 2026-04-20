import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.dispatch_notification import dispatch_notification


class DispatchNotificationTests(unittest.TestCase):
    def test_dispatches_email_notifications(self):
        result = dispatch_notification(
            "email",
            {
                "recipient": "dev@example.com",
                "subject": "Build finished",
                "message": "The deployment completed successfully.",
            },
        )

        self.assertEqual(
            {
                "recipient": "dev@example.com",
                "subject": "Build finished",
                "body": "The deployment completed successfully.",
            },
            result,
        )

    def test_uses_default_subject_for_email(self):
        result = dispatch_notification(
            "email",
            {
                "recipient": "dev@example.com",
                "message": "The deployment completed successfully.",
            },
        )

        self.assertEqual("Notification", result["subject"])

    def test_dispatches_sms_notifications(self):
        result = dispatch_notification(
            "sms",
            {
                "recipient": "+15550001111",
                "message": "Build finished successfully.",
            },
        )

        self.assertEqual(
            {
                "recipient": "+15550001111",
                "text": "Build finished successfully.",
            },
            result,
        )

    def test_dispatches_webhook_notifications(self):
        result = dispatch_notification(
            "webhook",
            {"url": "https://example.com/hook", "message": "deploy ok"},
        )

        self.assertEqual("https://example.com/hook", result["url"])
        self.assertEqual("POST", result["method"])
        self.assertEqual({"message": "deploy ok"}, result["json"])

    def test_rejects_unknown_channels(self):
        with self.assertRaises(ValueError):
            dispatch_notification("pagerduty", {"message": "test"})


if __name__ == "__main__":
    unittest.main()