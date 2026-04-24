import shutil
import socket
import subprocess
import tempfile
import time
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "update_rows_with_join.sql"


class MySQLSnippetTestCase(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.mysql = cls.require_binary("mysql")
        cls.mysqladmin = cls.require_binary("mysqladmin")
        cls.mysqld = cls.require_binary("mysqld")
        cls.mysql_share = cls.find_mysql_share()
        cls.temp_dir = tempfile.TemporaryDirectory(prefix="update-rows-join-mysql-")
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


class UpdateRowsWithJoinWithSqlTests(MySQLSnippetTestCase):
    def setUp(self):
        self.run_mysql(
            """
            DROP DATABASE IF EXISTS snippet_test;
            CREATE DATABASE snippet_test;
            USE snippet_test;

            CREATE TABLE items (
                item_id INT PRIMARY KEY,
                category_name VARCHAR(100) NOT NULL,
                display_rank INT NULL
            );

            CREATE TABLE category_defaults (
                category_name VARCHAR(100) PRIMARY KEY,
                display_rank INT NOT NULL
            );
            """
        )

    def insert_items(self, rows):
        values = ",\n".join(
            f"({item_id}, '{category_name}', {display_rank})"
            for item_id, category_name, display_rank in rows
        )
        self.run_mysql(
            f"""
            USE snippet_test;
            INSERT INTO items (item_id, category_name, display_rank)
            VALUES
            {values};
            """
        )

    def insert_defaults(self, rows):
        values = ",\n".join(
            f"('{category_name}', {display_rank})" for category_name, display_rank in rows
        )
        self.run_mysql(
            f"""
            USE snippet_test;
            INSERT INTO category_defaults (category_name, display_rank)
            VALUES
            {values};
            """
        )

    def run_update(self):
        self.run_mysql(f"USE snippet_test;\n{self.sql}")

    def fetch_items(self):
        output = self.run_mysql(
            """
            USE snippet_test;
            SELECT
                item_id,
                category_name,
                IFNULL(CAST(display_rank AS CHAR), 'NULL') AS display_rank
            FROM items
            ORDER BY item_id ASC;
            """,
            batch=True,
        )
        if not output:
            return []
        return [tuple(line.split("\t")) for line in output.splitlines()]

    def test_updates_matching_rows_with_missing_values(self):
        self.insert_items([(1, "alpha", "NULL"), (2, "beta", "NULL")])
        self.insert_defaults([("alpha", 10), ("beta", 20)])

        self.run_update()

        self.assertEqual(
            [("1", "alpha", "10"), ("2", "beta", "20")],
            self.fetch_items(),
        )

    def test_preserves_rows_that_already_have_a_value(self):
        self.insert_items([(1, "alpha", 99)])
        self.insert_defaults([("alpha", 10)])

        self.run_update()

        self.assertEqual([("1", "alpha", "99")], self.fetch_items())

    def test_leaves_unmatched_rows_unchanged(self):
        self.insert_items([(1, "gamma", "NULL")])
        self.insert_defaults([("alpha", 10)])

        self.run_update()

        self.assertEqual([("1", "gamma", "NULL")], self.fetch_items())

    def test_updates_multiple_rows_that_share_one_lookup_value(self):
        self.insert_items(
            [(1, "alpha", "NULL"), (2, "alpha", "NULL"), (3, "beta", "NULL")]
        )
        self.insert_defaults([("alpha", 10), ("beta", 20)])

        self.run_update()

        self.assertEqual(
            [("1", "alpha", "10"), ("2", "alpha", "10"), ("3", "beta", "20")],
            self.fetch_items(),
        )

    def test_does_nothing_when_the_lookup_table_is_empty(self):
        self.insert_items([(1, "alpha", "NULL"), (2, "beta", 7)])

        self.run_update()

        self.assertEqual(
            [("1", "alpha", "NULL"), ("2", "beta", "7")],
            self.fetch_items(),
        )


if __name__ == "__main__":
    unittest.main()