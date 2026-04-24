import sqlite3
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "calculate_running_totals.sql"


def read_sql() -> str:
    return SQL_PATH.read_text(encoding="utf-8")


class CalculateRunningTotalsWithSqlTests(unittest.TestCase):
    def setUp(self):
        self.connection = sqlite3.connect(":memory:")
        self.connection.row_factory = sqlite3.Row
        self.cursor = self.connection.cursor()
        self.cursor.execute(
            """
            CREATE TABLE ledger_entries (
                transaction_id INTEGER PRIMARY KEY,
                account_id TEXT NOT NULL,
                posted_at TEXT NOT NULL,
                amount INTEGER NOT NULL
            )
            """
        )
        self.sql = read_sql()

    def tearDown(self):
        self.connection.close()

    def insert_rows(self, rows):
        self.cursor.executemany(
            """
            INSERT INTO ledger_entries (transaction_id, account_id, posted_at, amount)
            VALUES (?, ?, ?, ?)
            """,
            rows,
        )
        self.connection.commit()

    def fetch_rows(self):
        rows = self.cursor.execute(self.sql).fetchall()
        return [dict(row) for row in rows]

    def test_returns_running_total_for_each_account(self):
        self.insert_rows(
            [
                (1, "acct-a", "2026-04-01T09:00:00", 50),
                (2, "acct-a", "2026-04-01T10:00:00", -20),
                (3, "acct-b", "2026-04-01T09:30:00", 100),
                (4, "acct-b", "2026-04-01T10:30:00", 25),
            ]
        )

        rows = self.fetch_rows()

        self.assertEqual(
            [
                {
                    "transaction_id": 1,
                    "account_id": "acct-a",
                    "posted_at": "2026-04-01T09:00:00",
                    "amount": 50,
                    "running_total": 50,
                },
                {
                    "transaction_id": 2,
                    "account_id": "acct-a",
                    "posted_at": "2026-04-01T10:00:00",
                    "amount": -20,
                    "running_total": 30,
                },
                {
                    "transaction_id": 3,
                    "account_id": "acct-b",
                    "posted_at": "2026-04-01T09:30:00",
                    "amount": 100,
                    "running_total": 100,
                },
                {
                    "transaction_id": 4,
                    "account_id": "acct-b",
                    "posted_at": "2026-04-01T10:30:00",
                    "amount": 25,
                    "running_total": 125,
                },
            ],
            rows,
        )

    def test_uses_transaction_id_as_tie_breaker_for_matching_timestamps(self):
        self.insert_rows(
            [
                (1, "acct-a", "2026-04-01T09:00:00", 40),
                (2, "acct-a", "2026-04-01T09:00:00", 10),
                (3, "acct-a", "2026-04-01T09:00:00", -5),
            ]
        )

        rows = self.fetch_rows()

        self.assertEqual([40, 50, 45], [row["running_total"] for row in rows])
        self.assertEqual([1, 2, 3], [row["transaction_id"] for row in rows])

    def test_handles_single_row_account(self):
        self.insert_rows([(1, "acct-a", "2026-04-01T09:00:00", 25)])

        rows = self.fetch_rows()

        self.assertEqual(
            [
                {
                    "transaction_id": 1,
                    "account_id": "acct-a",
                    "posted_at": "2026-04-01T09:00:00",
                    "amount": 25,
                    "running_total": 25,
                }
            ],
            rows,
        )

    def test_returns_empty_result_for_empty_table(self):
        self.assertEqual([], self.fetch_rows())

    def test_keeps_accounts_grouped_and_rows_ordered_inside_each_group(self):
        self.insert_rows(
            [
                (5, "acct-b", "2026-04-01T12:00:00", 10),
                (1, "acct-a", "2026-04-01T10:00:00", 5),
                (3, "acct-b", "2026-04-01T11:00:00", 30),
                (2, "acct-a", "2026-04-01T11:00:00", 15),
            ]
        )

        rows = self.fetch_rows()

        self.assertEqual(
            [
                ("acct-a", 1, 5),
                ("acct-a", 2, 20),
                ("acct-b", 3, 30),
                ("acct-b", 5, 40),
            ],
            [(row["account_id"], row["transaction_id"], row["running_total"]) for row in rows],
        )

    def test_supports_negative_amounts_in_running_total(self):
        self.insert_rows(
            [
                (1, "acct-a", "2026-04-01T09:00:00", 100),
                (2, "acct-a", "2026-04-01T10:00:00", -70),
                (3, "acct-a", "2026-04-01T11:00:00", -10),
            ]
        )

        rows = self.fetch_rows()

        self.assertEqual([100, 30, 20], [row["running_total"] for row in rows])


if __name__ == "__main__":
    unittest.main()