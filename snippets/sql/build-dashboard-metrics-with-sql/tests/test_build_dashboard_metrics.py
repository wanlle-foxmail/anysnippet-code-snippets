import sqlite3
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "build_dashboard_metrics.sql"


def read_sql() -> str:
    return SQL_PATH.read_text(encoding="utf-8")


class BuildDashboardMetricsWithSqlTests(unittest.TestCase):
    def setUp(self):
        self.connection = sqlite3.connect(":memory:")
        self.connection.row_factory = sqlite3.Row
        self.cursor = self.connection.cursor()
        self.cursor.execute(
            """
            CREATE TABLE orders (
                order_id INTEGER PRIMARY KEY AUTOINCREMENT,
                customer_id INTEGER NOT NULL,
                status TEXT NOT NULL,
                total_amount REAL,
                created_at TEXT NOT NULL,
                refunded_at TEXT
            )
            """
        )
        self.sql = read_sql()

    def tearDown(self):
        self.connection.close()

    def insert_orders(self, rows):
        self.cursor.executemany(
            """
            INSERT INTO orders (customer_id, status, total_amount, created_at, refunded_at)
            VALUES (?, ?, ?, ?, ?)
            """,
            rows,
        )
        self.connection.commit()

    def fetch_metrics(self, start_at: str, end_at: str):
        row = self.cursor.execute(
            self.sql,
            {"window_start": start_at, "window_end": end_at},
        ).fetchone()
        return dict(row)

    def test_counts_orders_by_status_in_single_query(self):
        self.insert_orders(
            [
                (101, "paid", 120.0, "2026-04-01T09:00:00", None),
                (102, "pending", 45.0, "2026-04-01T10:00:00", None),
                (103, "cancelled", 99.0, "2026-04-01T11:00:00", None),
                (104, "paid", 80.0, "2026-04-01T12:00:00", None),
            ]
        )

        metrics = self.fetch_metrics("2026-04-01T00:00:00", "2026-04-02T00:00:00")

        self.assertEqual(4, metrics["total_orders"])
        self.assertEqual(2, metrics["paid_orders"])
        self.assertEqual(1, metrics["pending_orders"])
        self.assertEqual(1, metrics["cancelled_orders"])

    def test_sums_paid_revenue_and_refunded_amount(self):
        self.insert_orders(
            [
                (101, "paid", 120.0, "2026-04-01T09:00:00", None),
                (102, "paid", 80.0, "2026-04-01T10:00:00", "2026-04-01T18:00:00"),
                (103, "pending", 45.0, "2026-04-01T11:00:00", None),
            ]
        )

        metrics = self.fetch_metrics("2026-04-01T00:00:00", "2026-04-02T00:00:00")

        self.assertEqual(200.0, metrics["paid_revenue"])
        self.assertEqual(1, metrics["refunded_orders"])
        self.assertEqual(80.0, metrics["refunded_amount"])

    def test_counts_refunds_by_refunded_at_window_instead_of_created_at(self):
        self.insert_orders(
            [
                (101, "paid", 70.0, "2026-03-30T09:00:00", "2026-04-01T08:00:00"),
                (102, "paid", 50.0, "2026-04-01T10:00:00", "2026-04-02T08:00:00"),
            ]
        )

        metrics = self.fetch_metrics("2026-04-01T00:00:00", "2026-04-02T00:00:00")

        self.assertEqual(1, metrics["total_orders"])
        self.assertEqual(1, metrics["paid_orders"])
        self.assertEqual(50.0, metrics["paid_revenue"])
        self.assertEqual(1, metrics["refunded_orders"])
        self.assertEqual(70.0, metrics["refunded_amount"])

    def test_returns_zeroed_metrics_for_empty_range(self):
        metrics = self.fetch_metrics("2026-04-01T00:00:00", "2026-04-02T00:00:00")

        self.assertEqual(
            {
                "total_orders": 0,
                "paid_orders": 0,
                "pending_orders": 0,
                "cancelled_orders": 0,
                "refunded_orders": 0,
                "paid_revenue": 0,
                "refunded_amount": 0,
            },
            metrics,
        )

    def test_treats_null_amount_as_zero_in_sums(self):
        self.insert_orders(
            [
                (101, "paid", None, "2026-04-01T09:00:00", None),
                (102, "paid", None, "2026-04-01T10:00:00", "2026-04-01T12:00:00"),
            ]
        )

        metrics = self.fetch_metrics("2026-04-01T00:00:00", "2026-04-02T00:00:00")

        self.assertEqual(2, metrics["paid_orders"])
        self.assertEqual(1, metrics["refunded_orders"])
        self.assertEqual(0, metrics["paid_revenue"])
        self.assertEqual(0, metrics["refunded_amount"])

    def test_uses_start_inclusive_end_exclusive_boundaries_for_created_and_refunded_at(self):
        self.insert_orders(
            [
                (101, "paid", 20.0, "2026-04-01T00:00:00", "2026-04-01T00:00:00"),
                (102, "paid", 30.0, "2026-04-02T00:00:00", "2026-04-02T00:00:00"),
                (103, "paid", 40.0, "2026-03-31T23:59:59", "2026-04-01T12:00:00"),
            ]
        )

        metrics = self.fetch_metrics("2026-04-01T00:00:00", "2026-04-02T00:00:00")

        self.assertEqual(1, metrics["total_orders"])
        self.assertEqual(20.0, metrics["paid_revenue"])
        self.assertEqual(2, metrics["refunded_orders"])
        self.assertEqual(60.0, metrics["refunded_amount"])


if __name__ == "__main__":
    unittest.main()