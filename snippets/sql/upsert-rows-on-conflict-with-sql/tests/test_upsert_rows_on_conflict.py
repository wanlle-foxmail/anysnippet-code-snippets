import sqlite3
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "upsert_rows_on_conflict.sql"


def read_sql() -> str:
    return SQL_PATH.read_text(encoding="utf-8")


class UpsertRowsOnConflictWithSqlTests(unittest.TestCase):
    def setUp(self):
        self.connection = sqlite3.connect(":memory:")
        self.connection.row_factory = sqlite3.Row
        self.cursor = self.connection.cursor()
        self.cursor.execute(
            """
            CREATE TABLE customer_contacts (
                contact_id INTEGER PRIMARY KEY AUTOINCREMENT,
                tenant_id INTEGER NOT NULL,
                email TEXT NOT NULL,
                full_name TEXT NOT NULL,
                phone TEXT,
                updated_at TEXT NOT NULL,
                UNIQUE (tenant_id, email)
            )
            """
        )
        self.sql = read_sql()

    def tearDown(self):
        self.connection.close()

    def insert_contact(self, row):
        self.cursor.execute(
            """
            INSERT INTO customer_contacts (tenant_id, email, full_name, phone, updated_at)
            VALUES (?, ?, ?, ?, ?)
            """,
            row,
        )
        self.connection.commit()
        return self.cursor.lastrowid

    def upsert_contact(self, row):
        self.cursor.execute(self.sql, row)
        self.connection.commit()

    def fetch_contacts(self):
        rows = self.cursor.execute(
            """
            SELECT contact_id, tenant_id, email, full_name, phone, updated_at
            FROM customer_contacts
            ORDER BY tenant_id, email
            """
        ).fetchall()
        return [dict(row) for row in rows]

    def test_inserts_new_contact_when_no_conflict_exists(self):
        self.upsert_contact((1, "dev@example.com", "Dev One", "+15550000001", "2026-04-19T10:00:00"))

        self.assertEqual(
            [
                {
                    "contact_id": 1,
                    "tenant_id": 1,
                    "email": "dev@example.com",
                    "full_name": "Dev One",
                    "phone": "+15550000001",
                    "updated_at": "2026-04-19T10:00:00",
                }
            ],
            self.fetch_contacts(),
        )

    def test_updates_existing_contact_when_unique_key_conflicts(self):
        contact_id = self.insert_contact(
            (1, "dev@example.com", "Dev One", "+15550000001", "2026-04-19T10:00:00")
        )

        self.upsert_contact((1, "dev@example.com", "Dev Renamed", "+15550009999", "2026-04-19T12:30:00"))

        contacts = self.fetch_contacts()
        self.assertEqual(1, len(contacts))
        self.assertEqual(contact_id, contacts[0]["contact_id"])
        self.assertEqual("Dev Renamed", contacts[0]["full_name"])
        self.assertEqual("+15550009999", contacts[0]["phone"])
        self.assertEqual("2026-04-19T12:30:00", contacts[0]["updated_at"])

    def test_keeps_same_email_separate_across_tenants(self):
        self.upsert_contact((1, "shared@example.com", "Tenant One", "+15550000001", "2026-04-19T10:00:00"))
        self.upsert_contact((2, "shared@example.com", "Tenant Two", "+15550000002", "2026-04-19T10:05:00"))

        contacts = self.fetch_contacts()
        self.assertEqual(2, len(contacts))
        self.assertEqual([1, 2], [contact["tenant_id"] for contact in contacts])
        self.assertEqual(["shared@example.com", "shared@example.com"], [contact["email"] for contact in contacts])

    def test_allows_null_phone_values_when_updating_conflicting_row(self):
        self.insert_contact((1, "nullable@example.com", "Dev One", "+15550000001", "2026-04-19T10:00:00"))

        self.upsert_contact((1, "nullable@example.com", "Dev One", None, "2026-04-19T11:00:00"))

        contacts = self.fetch_contacts()
        self.assertEqual(1, len(contacts))
        self.assertIsNone(contacts[0]["phone"])
        self.assertEqual("2026-04-19T11:00:00", contacts[0]["updated_at"])

    def test_supports_batch_upserts_with_executemany(self):
        self.insert_contact((1, "alex@example.com", "Alex Old", "+15550000001", "2026-04-19T09:00:00"))

        self.cursor.executemany(
            self.sql,
            [
                (1, "alex@example.com", "Alex New", "+15550001111", "2026-04-19T10:00:00"),
                (1, "bailey@example.com", "Bailey", "+15550002222", "2026-04-19T10:05:00"),
                (2, "alex@example.com", "Alex Tenant Two", "+15550003333", "2026-04-19T10:10:00"),
            ],
        )
        self.connection.commit()

        contacts = self.fetch_contacts()
        self.assertEqual(3, len(contacts))
        self.assertEqual("Alex New", contacts[0]["full_name"])
        self.assertEqual("Bailey", contacts[1]["full_name"])
        self.assertEqual("Alex Tenant Two", contacts[2]["full_name"])

    def test_does_not_change_unrelated_rows_when_upserting_conflicting_row(self):
        self.insert_contact((1, "keep@example.com", "Keep Me", "+15550000001", "2026-04-19T08:00:00"))
        self.insert_contact((1, "update@example.com", "Old Name", "+15550000002", "2026-04-19T08:30:00"))

        self.upsert_contact((1, "update@example.com", "New Name", "+15550009999", "2026-04-19T11:30:00"))

        contacts = self.fetch_contacts()
        self.assertEqual("Keep Me", contacts[0]["full_name"])
        self.assertEqual("+15550000001", contacts[0]["phone"])
        self.assertEqual("New Name", contacts[1]["full_name"])
        self.assertEqual("+15550009999", contacts[1]["phone"])


if __name__ == "__main__":
    unittest.main()