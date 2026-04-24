import csv
from pathlib import Path


def write_csv_part(output_dir: Path, header: list[str], rows: list[list[str]], part_number: int) -> str:
    output_path = output_dir / f"part-{part_number:03d}.csv"
    with open(output_path, "w", encoding="utf-8", newline="") as output_file:
        writer = csv.writer(output_file)
        writer.writerow(header)
        writer.writerows(rows)
    return str(output_path)


def split_csv_file_by_row_count(input_path: str, output_dir: str, rows_per_file: int) -> list[str]:
    """Split a CSV file into ordered parts with a fixed number of data rows per file."""
    # Flow:
    #   input CSV -> read header once
    #                |
    #                +-> empty file -> return []
    #                +-> collect rows until a part is full
    #                     |
    #                     +-> full part -> write it and reset the buffer
    #                     `-> leftover rows -> write the final part
    if isinstance(rows_per_file, bool) or not isinstance(rows_per_file, int) or rows_per_file <= 0:
        raise ValueError("rows_per_file must be a positive integer")

    output_dir_path = Path(output_dir)
    output_dir_path.mkdir(parents=True, exist_ok=True)

    with open(input_path, "r", encoding="utf-8", newline="") as input_file:
        reader = csv.reader(input_file)
        header = next(reader, None)
        if header is None:
            return []

        written_files = []
        pending_rows = []
        part_number = 1

        for row in reader:
            pending_rows.append(row)
            if len(pending_rows) == rows_per_file:
                written_files.append(write_csv_part(output_dir_path, header, pending_rows, part_number))
                pending_rows = []
                part_number += 1

        if pending_rows:
            written_files.append(write_csv_part(output_dir_path, header, pending_rows, part_number))

    return written_files


if __name__ == "__main__":
    result = split_csv_file_by_row_count("events.csv", "csv-parts", 1000)
    print(result)