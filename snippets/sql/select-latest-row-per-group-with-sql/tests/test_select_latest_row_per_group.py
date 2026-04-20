import sqlite3
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "select_latest_row_per_group.sql"


def read_sql() -> str:
    return SQL_PATH.read_text(encoding="utf-8")


class SelectLatestRowPerGroupWithSqlTests(unittest.TestCase):
    def setUp(self):
        self.connection = sqlite3.connect(":memory:")
        self.connection.row_factory = sqlite3.Row
        self.cursor = self.connection.cursor()
        self.cursor.execute(
            """
            CREATE TABLE device_status_events (
                event_id INTEGER PRIMARY KEY AUTOINCREMENT,
                device_id TEXT NOT NULL,
                status TEXT NOT NULL,
                recorded_at TEXT NOT NULL
            )
            """
        )
        self.sql = read_sql()

    def tearDown(self):
        self.connection.close()

    def insert_events(self, rows):
        self.cursor.executemany(
            """
            INSERT INTO device_status_events (device_id, status, recorded_at)
            VALUES (?, ?, ?)
            """,
            rows,
        )
        self.connection.commit()

    def fetch_latest_rows(self):
        rows = self.cursor.execute(self.sql).fetchall()
        return [dict(row) for row in rows]

    def test_returns_latest_event_for_each_device(self):
        self.insert_events(
            [
                ("device-a", "online", "2026-04-01T09:00:00"),
                ("device-a", "offline", "2026-04-01T11:00:00"),
                ("device-b", "online", "2026-04-01T10:00:00"),
                ("device-b", "sleeping", "2026-04-01T12:00:00"),
            ]
        )

        rows = self.fetch_latest_rows()

        self.assertEqual(2, len(rows))
        self.assertEqual("offline", rows[0]["status"])
        self.assertEqual("sleeping", rows[1]["status"])

    def test_uses_event_id_as_tie_breaker_for_matching_timestamps(self):
        self.insert_events(
            [
                ("device-a", "starting", "2026-04-01T09:00:00"),
                ("device-a", "online", "2026-04-01T09:00:00"),
            ]
        )

        rows = self.fetch_latest_rows()

        self.assertEqual(1, len(rows))
        self.assertEqual("online", rows[0]["status"])
        self.assertEqual(2, rows[0]["event_id"])

    def test_returns_single_row_for_device_with_one_event(self):
        self.insert_events([("device-a", "online", "2026-04-01T09:00:00")])

        rows = self.fetch_latest_rows()

        self.assertEqual(
            [
                {
                    "event_id": 1,
                    "device_id": "device-a",
                    "status": "online",
                    "recorded_at": "2026-04-01T09:00:00",
                }
            ],
            rows,
        )

    def test_returns_empty_result_for_empty_table(self):
        self.assertEqual([], self.fetch_latest_rows())

    def test_orders_result_by_device_id(self):
        self.insert_events(
            [
                ("device-c", "online", "2026-04-01T09:00:00"),
                ("device-a", "online", "2026-04-01T10:00:00"),
                ("device-b", "offline", "2026-04-01T11:00:00"),
            ]
        )

        rows = self.fetch_latest_rows()

        self.assertEqual(["device-a", "device-b", "device-c"], [row["device_id"] for row in rows])

    def test_ignores_older_rows_after_newer_row_exists(self):
        self.insert_events(
            [
                ("device-a", "online", "2026-04-01T08:00:00"),
                ("device-a", "sleeping", "2026-04-01T09:00:00"),
                ("device-a", "offline", "2026-04-01T10:00:00"),
            ]
        )

        rows = self.fetch_latest_rows()

        self.assertEqual(1, len(rows))
        self.assertEqual("offline", rows[0]["status"])
        self.assertEqual("2026-04-01T10:00:00", rows[0]["recorded_at"])


if __name__ == "__main__":
    unittest.main()