import hashlib
import io
import os
import sys
import tempfile
import unittest
from unittest.mock import MagicMock, patch

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "src"))

from download_large_file import download_large_file, stream_write_and_hash


class TestStreamWriteAndHash(unittest.TestCase):
    def test_writes_chunks_and_returns_correct_hash(self):
        chunks = [b"hello ", b"world"]
        buf = io.BytesIO()
        size, md5 = stream_write_and_hash(chunks, buf)
        self.assertEqual(size, 11)
        self.assertEqual(buf.getvalue(), b"hello world")
        expected_md5 = hashlib.md5(b"hello world").hexdigest()
        self.assertEqual(md5, expected_md5)


class TestDownloadLargeFile(unittest.TestCase):
    def _mock_response(self, chunks, status_code=200):
        resp = MagicMock()
        resp.status_code = status_code
        resp.raise_for_status = MagicMock()
        resp.iter_content = MagicMock(return_value=iter(chunks))
        resp.__enter__ = MagicMock(return_value=resp)
        resp.__exit__ = MagicMock(return_value=False)
        return resp

    @patch("download_large_file.requests.get")
    def test_downloads_multiple_chunks(self, mock_get):
        chunks = [b"aaa", b"bbb", b"ccc"]
        mock_get.return_value = self._mock_response(chunks)
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "out.bin")
            result = download_large_file("http://test.example/f", path)
            self.assertEqual(result["size"], 9)
            self.assertEqual(result["path"], path)
            expected_md5 = hashlib.md5(b"aaabbbccc").hexdigest()
            self.assertEqual(result["hash"], expected_md5)
            with open(path, "rb") as f:
                self.assertEqual(f.read(), b"aaabbbccc")

    @patch("download_large_file.requests.get")
    def test_empty_response_body(self, mock_get):
        mock_get.return_value = self._mock_response([])
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "empty.bin")
            result = download_large_file("http://test.example/f", path)
            self.assertEqual(result["size"], 0)
            expected_md5 = hashlib.md5(b"").hexdigest()
            self.assertEqual(result["hash"], expected_md5)

    @patch("download_large_file.requests.get")
    def test_http_error_cleans_up(self, mock_get):
        import requests as req

        resp = self._mock_response([])
        resp.raise_for_status.side_effect = req.exceptions.HTTPError("404")
        mock_get.return_value = resp
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "fail.bin")
            with self.assertRaises(req.exceptions.HTTPError):
                download_large_file("http://test.example/f", path)
            self.assertFalse(os.path.exists(path))

    @patch("download_large_file.requests.get")
    def test_timeout_cleans_up(self, mock_get):
        import requests as req

        mock_get.side_effect = req.exceptions.Timeout("timed out")
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "timeout.bin")
            with self.assertRaises(req.exceptions.Timeout):
                download_large_file("http://test.example/f", path)
            self.assertFalse(os.path.exists(path))

    @patch("download_large_file.requests.get")
    def test_mid_stream_disconnect_cleans_up(self, mock_get):
        import requests as req

        def failing_iter(chunk_size):
            yield b"partial"
            raise req.exceptions.ConnectionError("disconnected")

        resp = MagicMock()
        resp.raise_for_status = MagicMock()
        resp.iter_content = failing_iter
        resp.__enter__ = MagicMock(return_value=resp)
        resp.__exit__ = MagicMock(return_value=False)
        mock_get.return_value = resp
        with tempfile.TemporaryDirectory() as td:
            path = os.path.join(td, "partial.bin")
            with self.assertRaises(req.exceptions.ConnectionError):
                download_large_file("http://test.example/f", path)
            self.assertFalse(os.path.exists(path))


if __name__ == "__main__":
    unittest.main()
