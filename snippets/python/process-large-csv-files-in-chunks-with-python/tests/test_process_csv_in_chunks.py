import sys
import tempfile
import unittest
from pathlib import Path

import pandas as pd


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.process_csv_in_chunks import process_csv_in_chunks


class ProcessCsvInChunksTests(unittest.TestCase):
    def test_returns_single_processed_result_for_small_csv(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "sales.csv"
            pd.DataFrame(
                [
                    {"name": "A", "amount": 10},
                    {"name": "B", "amount": 15},
                ]
            ).to_csv(csv_path, index=False)

            result = process_csv_in_chunks(
                csv_path,
                lambda chunk: int(chunk["amount"].sum()),
                chunk_size=10,
            )

            self.assertEqual([25], result)

    def test_processes_multiple_chunks_in_read_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "rows.csv"
            pd.DataFrame(
                [
                    {"row_id": 1},
                    {"row_id": 2},
                    {"row_id": 3},
                    {"row_id": 4},
                    {"row_id": 5},
                ]
            ).to_csv(csv_path, index=False)

            result = process_csv_in_chunks(
                csv_path,
                lambda chunk: chunk["row_id"].tolist(),
                chunk_size=2,
            )

            self.assertEqual([[1, 2], [3, 4], [5]], result)

    def test_passes_read_csv_kwargs_to_pandas_reader(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "typed.csv"
            pd.DataFrame(
                [
                    {"name": "A", "amount": "10", "ignored": "x"},
                    {"name": "B", "amount": "15", "ignored": "y"},
                ]
            ).to_csv(csv_path, index=False)

            result = process_csv_in_chunks(
                csv_path,
                lambda chunk: (chunk.columns.tolist(), str(chunk["amount"].dtype)),
                chunk_size=1,
                read_csv_kwargs={"usecols": ["name", "amount"], "dtype": {"amount": "int64"}},
            )

            self.assertEqual(
                [(["name", "amount"], "int64"), (["name", "amount"], "int64")],
                result,
            )

    def test_returns_empty_result_for_header_only_csv(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "header-only.csv"
            csv_path.write_text("name,amount\n", encoding="utf-8")
            call_count = 0

            def processor(chunk):
                nonlocal call_count
                call_count += 1
                return chunk.shape[0]

            result = process_csv_in_chunks(csv_path, processor, chunk_size=2)

            self.assertEqual([], result)
            self.assertEqual(0, call_count)

    def test_raises_file_not_found_for_missing_csv(self):
        with self.assertRaises(FileNotFoundError):
            process_csv_in_chunks("missing.csv", lambda chunk: chunk.shape[0])

    def test_raises_value_error_for_invalid_chunk_size(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "data.csv"
            pd.DataFrame([{"value": 1}]).to_csv(csv_path, index=False)

            with self.assertRaises(ValueError):
                process_csv_in_chunks(csv_path, lambda chunk: chunk.shape[0], chunk_size=0)

            with self.assertRaises(ValueError):
                process_csv_in_chunks(csv_path, lambda chunk: chunk.shape[0], chunk_size=-1)

            with self.assertRaises(ValueError):
                process_csv_in_chunks(csv_path, lambda chunk: chunk.shape[0], chunk_size="2")

            with self.assertRaises(ValueError):
                process_csv_in_chunks(csv_path, lambda chunk: chunk.shape[0], chunk_size=True)

    def test_raises_value_error_for_non_callable_chunk_processor(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "data.csv"
            pd.DataFrame([{"value": 1}]).to_csv(csv_path, index=False)

            with self.assertRaises(ValueError):
                process_csv_in_chunks(csv_path, None, chunk_size=2)

    def test_propagates_chunk_processor_exceptions(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "data.csv"
            pd.DataFrame([{"value": 1}, {"value": 2}]).to_csv(csv_path, index=False)

            def processor(chunk):
                raise RuntimeError(f"bad chunk size: {len(chunk)}")

            with self.assertRaisesRegex(RuntimeError, "bad chunk size"):
                process_csv_in_chunks(csv_path, processor, chunk_size=1)

    def test_rejects_conflicting_chunksize_in_read_csv_kwargs(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "data.csv"
            pd.DataFrame([{"value": 1}]).to_csv(csv_path, index=False)

            with self.assertRaises(ValueError):
                process_csv_in_chunks(
                    csv_path,
                    lambda chunk: chunk.shape[0],
                    chunk_size=2,
                    read_csv_kwargs={"chunksize": 1},
                )


if __name__ == "__main__":
    unittest.main()