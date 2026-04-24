import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.compress_text_for_storage import compress_text, decompress_text


class CompressTextForStorageTests(unittest.TestCase):
    def test_round_trip_for_json_text(self):
        text = '{"name": "AnySnippet", "items": [1, 2, 3], "active": true}'

        compressed_text = compress_text(text)

        self.assertEqual(text, decompress_text(compressed_text))

    def test_round_trip_for_empty_string(self):
        self.assertEqual("", decompress_text(compress_text("")))

    def test_round_trip_preserves_unicode_text(self):
        text = "你好, AnySnippet. Compression should keep emoji 🚀 and accents cafe."

        compressed_text = compress_text(text)

        self.assertEqual(text, decompress_text(compressed_text))

    def test_reduces_size_for_repetitive_json_text(self):
        text = '{"events": [' + ','.join(['{"type": "view", "page": "/docs", "status": "ok"}'] * 80) + ']}'

        compressed_text = compress_text(text)

        self.assertLess(len(compressed_text), len(text))

    def test_decompress_raises_value_error_for_invalid_base64(self):
        with self.assertRaises(ValueError):
            decompress_text("not-base64")

    def test_decompress_raises_value_error_for_non_gzip_payload(self):
        with self.assertRaises(ValueError):
            decompress_text("cGxhaW4tdGV4dA==")

    def test_compress_rejects_non_string_input(self):
        with self.assertRaises(TypeError):
            compress_text(None)

        with self.assertRaises(TypeError):
            compress_text(123)

    def test_decompress_rejects_non_string_input(self):
        with self.assertRaises(TypeError):
            decompress_text(None)

        with self.assertRaises(TypeError):
            decompress_text(123)


if __name__ == "__main__":
    unittest.main()