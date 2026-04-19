import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.dispatch_notification import dispatch_notification


class DispatchNotificationTests(unittest.TestCase):
    def test_dispatches_email_notifications_with_default_handler(self):
        result = dispatch_notification(
            "email",
            {
                "recipient": "dev@example.com",
                "subject": "Build finished",
                "message": "The deployment completed successfully.",
            },
        )

        self.assertEqual("email", result["channel"])
        self.assertEqual("email", result["handled_by"])
        self.assertEqual(
            {
                "recipient": "dev@example.com",
                "subject": "Build finished",
                "body": "The deployment completed successfully.",
            },
            result["output"],
        )

    def test_dispatches_sms_notifications_with_channel_specific_logic(self):
        result = dispatch_notification(
            "sms",
            {
                "recipient": "+15550001111",
                "message": "  Build finished successfully.  ",
            },
        )

        self.assertEqual("sms", result["channel"])
        self.assertEqual(
            {
                "recipient": "+15550001111",
                "text": "Build finished successfully.",
            },
            result["output"],
        )

    def test_normalizes_channel_name_before_dispatch(self):
        result = dispatch_notification(
            "  WEBHOOK  ",
            {"url": "https://example.com/hook", "message": "deploy ok"},
        )

        self.assertEqual("webhook", result["channel"])
        self.assertEqual("webhook", result["handled_by"])
        self.assertEqual("https://example.com/hook", result["output"]["url"])

    def test_rejects_unknown_channels(self):
        with self.assertRaises(ValueError):
            dispatch_notification("pagerduty", {"message": "test"})

    def test_allows_custom_handler_registry(self):
        def slack_handler(payload):
            return {"channel": payload["channel"], "text": payload["message"].upper()}

        result = dispatch_notification(
            "slack",
            {"channel": "deployments", "message": "done"},
            handlers={" SLACK ": slack_handler},
        )

        self.assertEqual("slack", result["handled_by"])
        self.assertEqual({"channel": "deployments", "text": "DONE"}, result["output"])

    def test_passes_payload_copy_to_handler(self):
        original_payload = {"channel": "alerts", "message": "test"}

        def mutating_handler(payload):
            payload["message"] = "changed"
            return {"message": payload["message"]}

        result = dispatch_notification("slack", original_payload, handlers={"slack": mutating_handler})

        self.assertEqual({"channel": "alerts", "message": "test"}, original_payload)
        self.assertEqual({"message": "changed"}, result["output"])

    def test_rejects_blank_channel_names(self):
        with self.assertRaises(ValueError):
            dispatch_notification("   ", {"message": "test"})

        with self.assertRaises(ValueError):
            dispatch_notification(None, {"message": "test"})


if __name__ == "__main__":
    unittest.main()