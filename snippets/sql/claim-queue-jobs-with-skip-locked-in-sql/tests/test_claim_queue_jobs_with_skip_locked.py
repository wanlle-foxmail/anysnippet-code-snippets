import shutil
import socket
import subprocess
import tempfile
import time
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
SQL_PATH = SNIPPET_ROOT / "src" / "claim_queue_jobs_with_skip_locked.sql"


class MySQLSnippetTestCase(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.mysql = cls.require_binary("mysql")
        cls.mysqladmin = cls.require_binary("mysqladmin")
        cls.mysqld = cls.require_binary("mysqld")
        cls.mysql_share = cls.find_mysql_share()
        cls.temp_dir = tempfile.TemporaryDirectory(prefix="claim-jobs-mysql-")
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

    @classmethod
    def start_locking_transaction(cls, sql: str) -> subprocess.Popen:
        process = subprocess.Popen(
            [
                cls.mysql,
                "--protocol=TCP",
                "-h127.0.0.1",
                f"-P{cls.port}",
                "-uroot",
                "--batch",
                "--raw",
                "--skip-column-names",
                "--execute",
                sql,
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
        time.sleep(0.2)
        return process


class ClaimQueueJobsWithSkipLockedTests(MySQLSnippetTestCase):
    def setUp(self):
        self.run_mysql(
            """
            DROP DATABASE IF EXISTS snippet_test;
            CREATE DATABASE snippet_test;
            USE snippet_test;

            CREATE TABLE jobs (
                job_id INT PRIMARY KEY,
                payload VARCHAR(100) NOT NULL,
                status VARCHAR(20) NOT NULL,
                priority INT NOT NULL,
                worker_id VARCHAR(50) NULL,
                claimed_at DATETIME NULL
            );

            CREATE INDEX idx_jobs_status_priority_job_id
            ON jobs (status, priority DESC, job_id ASC);
            """
        )

    def insert_jobs(self, rows):
        values = ",\n".join(
            f"({job_id}, '{payload}', '{status}', {priority}, {worker_id}, {claimed_at})"
            for job_id, payload, status, priority, worker_id, claimed_at in rows
        )
        self.run_mysql(
            f"""
            USE snippet_test;
            INSERT INTO jobs (job_id, payload, status, priority, worker_id, claimed_at)
            VALUES
            {values};
            """
        )

    def run_snippet(self):
        output = self.run_mysql(f"USE snippet_test;\n{self.sql}", batch=True)
        if not output:
            return []
        return [tuple(line.split("\t")) for line in output.splitlines()]

    def fetch_job_states(self):
        output = self.run_mysql(
            """
            USE snippet_test;
            SELECT job_id, status, COALESCE(worker_id, ''), COALESCE(claimed_at, '')
            FROM jobs
            ORDER BY job_id ASC;
            """,
            batch=True,
        )
        if not output:
            return []
        return [tuple(line.split("\t")) for line in output.splitlines()]

    def test_claims_the_highest_priority_queued_jobs(self):
        self.insert_jobs(
            [
                (1, "job-a", "queued", 10, "NULL", "NULL"),
                (2, "job-b", "queued", 6, "NULL", "NULL"),
                (3, "job-c", "processing", 8, "'worker-x'", "'2026-04-01 08:30:00'"),
                (4, "job-d", "queued", 9, "NULL", "NULL"),
            ]
        )

        self.assertEqual(
            [("1", "10", "worker-1"), ("4", "9", "worker-1")],
            self.run_snippet(),
        )

    def test_skips_rows_locked_by_another_transaction(self):
        self.insert_jobs(
            [
                (1, "job-a", "queued", 10, "NULL", "NULL"),
                (2, "job-b", "queued", 9, "NULL", "NULL"),
                (3, "job-c", "queued", 8, "NULL", "NULL"),
            ]
        )

        locker = self.start_locking_transaction(
            "USE snippet_test; START TRANSACTION; SELECT job_id FROM jobs WHERE job_id = 1 FOR UPDATE; DO SLEEP(2); COMMIT;"
        )
        try:
            self.assertEqual(
                [("2", "9", "worker-1"), ("3", "8", "worker-1")],
                self.run_snippet(),
            )
        finally:
            locker.communicate(timeout=5)

    def test_limits_the_claim_batch_to_two_jobs(self):
        self.insert_jobs(
            [
                (1, "job-a", "queued", 10, "NULL", "NULL"),
                (2, "job-b", "queued", 9, "NULL", "NULL"),
                (3, "job-c", "queued", 8, "NULL", "NULL"),
                (4, "job-d", "queued", 7, "NULL", "NULL"),
            ]
        )

        self.assertEqual(2, len(self.run_snippet()))

    def test_breaks_priority_ties_by_job_id(self):
        self.insert_jobs(
            [
                (1, "job-a", "queued", 10, "NULL", "NULL"),
                (2, "job-b", "queued", 10, "NULL", "NULL"),
                (3, "job-c", "queued", 10, "NULL", "NULL"),
            ]
        )

        self.assertEqual(
            [("1", "10", "worker-1"), ("2", "10", "worker-1")],
            self.run_snippet(),
        )

    def test_ignores_non_queued_jobs(self):
        self.insert_jobs(
            [
                (1, "job-a", "done", 10, "'worker-a'", "'2026-04-01 08:00:00'"),
                (2, "job-b", "processing", 9, "'worker-b'", "'2026-04-01 08:15:00'"),
                (3, "job-c", "queued", 8, "NULL", "NULL"),
                (4, "job-d", "queued", 7, "NULL", "NULL"),
            ]
        )

        self.assertEqual(
            [("3", "8", "worker-1"), ("4", "7", "worker-1")],
            self.run_snippet(),
        )
        self.assertEqual(
            [
                ("1", "done", "worker-a", "2026-04-01 08:00:00"),
                ("2", "processing", "worker-b", "2026-04-01 08:15:00"),
                ("3", "processing", "worker-1", ""),
                ("4", "processing", "worker-1", ""),
            ],
            [(job_id, status, worker_id, "") if worker_id == "worker-1" else (job_id, status, worker_id, claimed_at) for job_id, status, worker_id, claimed_at in self.fetch_job_states()],
        )


if __name__ == "__main__":
    unittest.main()