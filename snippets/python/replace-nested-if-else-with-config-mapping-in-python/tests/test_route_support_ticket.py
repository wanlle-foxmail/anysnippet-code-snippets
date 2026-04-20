import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.route_support_ticket import route_support_ticket


class RouteSupportTicketTests(unittest.TestCase):
    def test_routes_enterprise_billing_tickets(self):
        result = route_support_ticket("chat", "enterprise", "billing")

        self.assertEqual("priority-billing", result["queue"])
        self.assertEqual("urgent", result["priority"])
        self.assertEqual(1, result["sla_hours"])

    def test_routes_pro_login_tickets(self):
        result = route_support_ticket("email", "pro", "login")

        self.assertEqual("pro-login", result["queue"])
        self.assertEqual("high", result["priority"])
        self.assertEqual(4, result["sla_hours"])

    def test_routes_enterprise_security_phone_tickets(self):
        result = route_support_ticket("phone", "enterprise", "security")

        self.assertEqual("security-incident", result["queue"])
        self.assertEqual("urgent", result["priority"])
        self.assertEqual(1, result["sla_hours"])

    def test_falls_back_to_default_rule(self):
        result = route_support_ticket("forum", "starter", "feature-request")

        self.assertEqual("general-support", result["queue"])
        self.assertEqual("normal", result["priority"])
        self.assertEqual(24, result["sla_hours"])

    def test_returns_copy_of_rule_data(self):
        first_result = route_support_ticket("chat", "enterprise", "billing")
        first_result["queue"] = "mutated"

        second_result = route_support_ticket("chat", "enterprise", "billing")

        self.assertEqual("priority-billing", second_result["queue"])


if __name__ == "__main__":
    unittest.main()