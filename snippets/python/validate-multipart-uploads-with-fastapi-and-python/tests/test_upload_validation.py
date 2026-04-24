import os
import sys
import unittest

from fastapi.testclient import TestClient


SNIPPET_ROOT = os.path.join(os.path.dirname(__file__), "..", "src")
if SNIPPET_ROOT not in sys.path:
    sys.path.insert(0, SNIPPET_ROOT)

from upload_validation import ALLOWED_CONTENT_TYPES, MAX_FILE_SIZE_BYTES, app


class UploadValidationTests(unittest.TestCase):
    def setUp(self):
        self.client = TestClient(app)

    def test_accepts_valid_text_upload(self):
        response = self.client.post(
            "/upload",
            files={"file": ("notes.txt", b"hello", "text/plain")},
        )

        self.assertEqual(200, response.status_code)
        self.assertEqual(
            {
                "filename": "notes.txt",
                "content_type": "text/plain",
                "size": 5,
            },
            response.json(),
        )

    def test_rejects_blank_filename(self):
        response = self.client.post(
            "/upload",
            files={"file": ("   ", b"hello", "text/plain")},
        )

        self.assertEqual(400, response.status_code)
        self.assertEqual("filename is required", response.json()["detail"])

    def test_rejects_unsupported_content_type(self):
        response = self.client.post(
            "/upload",
            files={"file": ("data.json", b"{}", "application/json")},
        )

        self.assertEqual(415, response.status_code)
        self.assertEqual("unsupported content type", response.json()["detail"])

    def test_rejects_empty_files(self):
        response = self.client.post(
            "/upload",
            files={"file": ("empty.txt", b"", "text/plain")},
        )

        self.assertEqual(400, response.status_code)
        self.assertEqual("file is empty", response.json()["detail"])

    def test_rejects_files_larger_than_the_limit(self):
        response = self.client.post(
            "/upload",
            files={"file": ("big.txt", b"a" * (MAX_FILE_SIZE_BYTES + 1), "text/plain")},
        )

        self.assertEqual(413, response.status_code)
        self.assertEqual("file is too large", response.json()["detail"])

    def test_accepts_files_exactly_at_the_size_limit(self):
        response = self.client.post(
            "/upload",
            files={"file": ("limit.csv", b"a" * MAX_FILE_SIZE_BYTES, "text/csv")},
        )

        self.assertEqual(200, response.status_code)
        self.assertEqual("limit.csv", response.json()["filename"])
        self.assertEqual("text/csv", response.json()["content_type"])
        self.assertEqual(MAX_FILE_SIZE_BYTES, response.json()["size"])

    def test_returns_a_sanitized_basename_for_uploaded_filenames(self):
        response = self.client.post(
            "/upload",
            files={"file": ("../reports/quarterly.csv", b"name\nvalue\n", "text/csv")},
        )

        self.assertEqual(200, response.status_code)
        self.assertEqual("quarterly.csv", response.json()["filename"])

    def test_allowed_content_types_stay_small_and_explicit(self):
        self.assertEqual({"text/plain", "text/csv"}, ALLOWED_CONTENT_TYPES)


if __name__ == "__main__":
    unittest.main()