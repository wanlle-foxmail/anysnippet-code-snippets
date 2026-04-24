import shutil
import socket
import subprocess
import tempfile
import time
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "find_missing_related_rows_with_left_join.sql"


class MySQLSnippetTestCase(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.mysql = cls.require_binary("mysql")
        cls.mysqladmin = cls.require_binary("mysqladmin")
        cls.mysqld = cls.require_binary("mysqld")
        cls.mysql_share = cls.find_mysql_share()
        cls.temp_dir = tempfile.TemporaryDirectory(prefix="missing-left-join-mysql-")
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


class FindMissingRelatedRowsWithLeftJoinWithSqlTests(MySQLSnippetTestCase):
    def setUp(self):
        self.run_mysql(
            """
            DROP DATABASE IF EXISTS snippet_test;
            CREATE DATABASE snippet_test;
            USE snippet_test;

            CREATE TABLE users (
                user_id INT PRIMARY KEY,
                username VARCHAR(100) NOT NULL
            );

            CREATE TABLE profiles (
                profile_id INT PRIMARY KEY,
                user_id INT NULL,
                profile_name VARCHAR(100) NOT NULL
            );
            """
        )

    def insert_users(self, rows):
        values = ",\n".join(f"({user_id}, '{username}')" for user_id, username in rows)
        self.run_mysql(
            f"""
            USE snippet_test;
            INSERT INTO users (user_id, username)
            VALUES
            {values};
            """
        )

    def insert_profiles(self, rows):
        values = ",\n".join(
            f"({profile_id}, {user_id}, '{profile_name}')"
            for profile_id, user_id, profile_name in rows
        )
        self.run_mysql(
            f"""
            USE snippet_test;
            INSERT INTO profiles (profile_id, user_id, profile_name)
            VALUES
            {values};
            """
        )

    def fetch_rows(self):
        output = self.run_mysql(f"USE snippet_test;\n{self.sql}", batch=True)
        if not output:
            return []
        return [tuple(line.split("\t")) for line in output.splitlines()]

    def test_finds_rows_that_have_no_related_match(self):
        self.insert_users([(1, "alpha"), (2, "bravo"), (3, "charlie")])
        self.insert_profiles([(1, 1, "profile-a"), (2, 3, "profile-c")])

        self.assertEqual([("2", "bravo")], self.fetch_rows())

    def test_returns_empty_result_when_every_row_is_matched(self):
        self.insert_users([(1, "alpha"), (2, "bravo")])
        self.insert_profiles([(1, 1, "profile-a"), (2, 2, "profile-b")])

        self.assertEqual([], self.fetch_rows())

    def test_returns_all_rows_when_the_related_table_is_empty(self):
        self.insert_users([(1, "alpha"), (2, "bravo")])

        self.assertEqual([("1", "alpha"), ("2", "bravo")], self.fetch_rows())

    def test_ignores_null_keys_on_the_right_table(self):
        self.insert_users([(1, "alpha"), (2, "bravo")])
        self.insert_profiles([(1, "NULL", "orphan-profile")])

        self.assertEqual([("1", "alpha"), ("2", "bravo")], self.fetch_rows())

    def test_orders_missing_rows_by_the_left_table_key(self):
        self.insert_users([(3, "charlie"), (1, "alpha"), (2, "bravo")])
        self.insert_profiles([(1, 2, "profile-b")])

        self.assertEqual([("1", "alpha"), ("3", "charlie")], self.fetch_rows())


if __name__ == "__main__":
    unittest.main()