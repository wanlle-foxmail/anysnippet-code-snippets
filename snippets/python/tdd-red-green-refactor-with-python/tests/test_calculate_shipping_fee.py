import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.calculate_shipping_fee import calculate_shipping_fee


class CalculateShippingFeeTests(unittest.TestCase):
    def test_uses_standard_tier_by_default_below_threshold(self):
        self.assertEqual(500, calculate_shipping_fee(3000))

    def test_returns_free_shipping_at_threshold(self):
        self.assertEqual(0, calculate_shipping_fee(5000, "standard"))

    def test_premium_customer_always_gets_free_shipping(self):
        self.assertEqual(0, calculate_shipping_fee(1000, "premium"))

    def test_rejects_negative_subtotal(self):
        with self.assertRaises(ValueError):
            calculate_shipping_fee(-100, "standard")

    def test_rejects_bool_subtotal(self):
        with self.assertRaises(TypeError):
            calculate_shipping_fee(True, "standard")

    def test_rejects_unknown_customer_tier(self):
        with self.assertRaises(ValueError):
            calculate_shipping_fee(3000, "gold")

    def test_charges_shipping_just_below_threshold(self):
        self.assertEqual(500, calculate_shipping_fee(4999, "standard"))


if __name__ == "__main__":
    unittest.main()