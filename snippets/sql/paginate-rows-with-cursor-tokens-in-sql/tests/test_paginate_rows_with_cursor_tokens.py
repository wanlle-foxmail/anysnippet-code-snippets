import shutil
import socket
import subprocess
import tempfile
import time
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "paginate_rows_with_cursor_tokens.sql"


class MySQLSnippetTestCase(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.mysql = cls.require_binary("mysql")
        cls.mysqladmin = cls.require_binary("mysqladmin")
        cls.mysqld = cls.require_binary("mysqld")
        cls.mysql_share = cls.find_mysql_share()
        cls.temp_dir = tempfile.TemporaryDirectory(prefix="paginate-cursor-mysql-")
        cls.server_root = Path(cls.temp_dir.name)
        cls.data_dir = cls.server_root / "data"
        cls.run_dir = cls.server_root / "run"
        cls.log_path = cls.server_root / "mysql.log"
        cls.socket_path = cls.run_dir / "mysql.sock"
        cls.pid_path = cls.run_dir / "mysql.pid"
        cls.port = cls.find_free_port()
        cls.data_dir.mkdir(parents=True, exist_ok=True)
        cls.run_dir.mkdir(parents=True, exist_ok=True)
        cls.initialize_server()
        cls.start_server()
        cls.sql = SQL_PATH.read_text(encoding="utf-8")

    @classmethod
    def tearDownClass(cls):
        try:
            cls.stop_server()
        finally:
            cls.temp_dir.cleanup()

    @classmethod
    def require_binary(cls, name: str) -> str:
        path = shutil.which(name)
        if not path:
            raise unittest.SkipTest(f"{name} is required for MySQL verification")
        return path

    @classmethod
    def find_mysql_share(cls):
        mysqld_path = Path(cls.mysqld).resolve()
        candidates = [
            mysqld_path.parent.parent / "share" / "mysql",
            Path("/usr/share/mysql"),
            Path("/usr/local/share/mysql"),
        ]
        package_root = mysqld_path.parent.parent / "pkgs"
        if package_root.is_dir():
            candidates.extend(sorted(package_root.glob("mysql-*/share/mysql")))

        for candidate in candidates:
            if (candidate / "english" / "errmsg.sys").exists():
                return str(candidate)
        return None

    @classmethod
    def find_free_port(cls) -> int:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            sock.bind(("127.0.0.1", 0))
            return sock.getsockname()[1]

    @classmethod
    def initialize_server(cls):
        command = [
            cls.mysqld,
            "--no-defaults",
            "--initialize-insecure",
            f"--datadir={cls.data_dir}",
        ]
        if cls.mysql_share:
            command.append(f"--lc-messages-dir={cls.mysql_share}")

        subprocess.run(command, capture_output=True, text=True, check=True)

    @classmethod
    def start_server(cls):
        command = [
            cls.mysqld,
            "--no-defaults",
            f"--datadir={cls.data_dir}",
            f"--socket={cls.socket_path}",
            f"--port={cls.port}",
            f"--pid-file={cls.pid_path}",
            "--bind-address=127.0.0.1",
            "--skip-networking=0",
            "--daemonize",
            "--mysqlx=0",
            f"--log-error={cls.log_path}",
        ]
        if cls.mysql_share:
            command.append(f"--lc-messages-dir={cls.mysql_share}")

        subprocess.run(command, capture_output=True, text=True, check=True)

        for _ in range(60):
            ping = subprocess.run(
                [
                    cls.mysqladmin,
                    "--protocol=TCP",
                    "-h127.0.0.1",
                    f"-P{cls.port}",
                    "-uroot",
                    "ping",
                    "--silent",
                ],
                capture_output=True,
                text=True,
            )
            if ping.returncode == 0:
                return
            time.sleep(0.1)

        raise AssertionError(cls.log_path.read_text(encoding="utf-8"))

    @classmethod
    def stop_server(cls):
        subprocess.run(
            [
                cls.mysqladmin,
                "--protocol=TCP",
                "-h127.0.0.1",
                f"-P{cls.port}",
                "-uroot",
                "shutdown",
            ],
            capture_output=True,
            text=True,
            check=True,
        )

    @classmethod
    def run_mysql(cls, sql: str, *, batch: bool = False) -> str:
        command = [
            cls.mysql,
            "--protocol=TCP",
            "-h127.0.0.1",
            f"-P{cls.port}",
            "-uroot",
        ]
        if batch:
            command.extend(["--batch", "--raw", "--skip-column-names"])

        completed = subprocess.run(
            command,
            input=sql,
            capture_output=True,
            text=True,
            check=True,
        )
        return completed.stdout.strip()


class PaginateRowsWithCursorTokensTests(MySQLSnippetTestCase):
    def setUp(self):
        self.run_mysql(
            """
            DROP DATABASE IF EXISTS snippet_test;
            CREATE DATABASE snippet_test;
            USE snippet_test;

            CREATE TABLE items (
                item_id INT PRIMARY KEY,
                item_name VARCHAR(100) NOT NULL,
                created_at DATETIME NOT NULL
            );

            CREATE INDEX idx_items_created_at_item_id
            ON items (created_at DESC, item_id DESC);
            """
        )

    def insert_items(self, rows):
        values = ",\n".join(
            f"({item_id}, '{item_name}', '{created_at}')"
            for item_id, item_name, created_at in rows
        )
        self.run_mysql(
            f"""
            USE snippet_test;
            INSERT INTO items (item_id, item_name, created_at)
            VALUES
            {values};
            """
        )

    def fetch_rows(self):
        output = self.run_mysql(f"USE snippet_test;\n{self.sql}", batch=True)
        if not output:
            return []
        return [tuple(line.split("\t")) for line in output.splitlines()]

    def test_returns_the_next_page_after_the_cursor(self):
        self.insert_items(
            [
                (1, "alpha", "2026-04-01 08:00:00"),
                (2, "bravo", "2026-04-01 09:00:00"),
                (3, "charlie", "2026-04-01 10:00:00"),
                (4, "delta", "2026-04-01 11:00:00"),
                (5, "echo", "2026-04-01 12:00:00"),
                (6, "foxtrot", "2026-04-01 13:00:00"),
                (7, "golf", "2026-04-01 14:00:00"),
            ]
        )

        self.assertEqual(
            [
                ("4", "delta", "2026-04-01 11:00:00"),
                ("3", "charlie", "2026-04-01 10:00:00"),
                ("2", "bravo", "2026-04-01 09:00:00"),
            ],
            self.fetch_rows(),
        )

    def test_limits_the_page_to_three_rows(self):
        self.insert_items(
            [
                (1, "alpha", "2026-04-01 08:00:00"),
                (2, "bravo", "2026-04-01 09:00:00"),
                (3, "charlie", "2026-04-01 10:00:00"),
                (4, "delta", "2026-04-01 11:00:00"),
                (5, "echo", "2026-04-01 12:00:00"),
                (6, "foxtrot", "2026-04-01 13:00:00"),
                (7, "golf", "2026-04-01 14:00:00"),
                (8, "hotel", "2026-04-01 15:00:00"),
            ]
        )

        self.assertEqual(3, len(self.fetch_rows()))

    def test_uses_item_id_as_a_cursor_tie_breaker(self):
        self.insert_items(
            [
                (1, "alpha", "2026-04-01 08:00:00"),
                (2, "bravo", "2026-04-01 09:00:00"),
                (3, "charlie", "2026-04-01 10:00:00"),
                (4, "delta", "2026-04-01 11:00:00"),
                (5, "echo", "2026-04-01 12:00:00"),
                (6, "foxtrot", "2026-04-01 12:00:00"),
                (7, "golf", "2026-04-01 13:00:00"),
            ]
        )

        self.assertEqual(
            [
                ("4", "delta", "2026-04-01 11:00:00"),
                ("3", "charlie", "2026-04-01 10:00:00"),
                ("2", "bravo", "2026-04-01 09:00:00"),
            ],
            self.fetch_rows(),
        )

    def test_returns_a_partial_page_near_the_end(self):
        self.insert_items(
            [
                (1, "alpha", "2026-04-01 08:00:00"),
                (2, "bravo", "2026-04-01 09:00:00"),
                (3, "charlie", "2026-04-01 10:00:00"),
                (4, "delta", "2026-04-01 11:00:00"),
                (5, "echo", "2026-04-01 12:00:00"),
            ]
        )

        self.assertEqual(
            [
                ("4", "delta", "2026-04-01 11:00:00"),
                ("3", "charlie", "2026-04-01 10:00:00"),
                ("2", "bravo", "2026-04-01 09:00:00"),
            ],
            self.fetch_rows(),
        )

    def test_returns_no_rows_when_the_cursor_is_already_at_the_end(self):
        self.insert_items(
            [
                (1, "alpha", "2026-04-01 08:00:00"),
                (2, "bravo", "2026-04-01 09:00:00"),
                (3, "charlie", "2026-04-01 10:00:00"),
                (4, "delta", "2026-04-01 11:00:00"),
                (5, "echo", "2026-04-01 12:00:00"),
            ]
        )
        self.run_mysql(
            """
            USE snippet_test;
            UPDATE items
            SET created_at = '2026-04-01 12:00:00'
            WHERE item_id = 5;
            """
        )

        self.run_mysql(
            """
            USE snippet_test;
            DELETE FROM items WHERE item_id IN (1, 2, 3, 4);
            """
        )

        self.assertEqual([], self.fetch_rows())


if __name__ == "__main__":
    unittest.main()