import json


def deduplicate_jsonl_records_by_key(input_path: str, output_path: str, key_name: str) -> int:
    """Write the first JSONL object for each unique top-level string key and return the count."""
    # Flow:
    #   input JSONL -> parse each non-empty line
    #                  |
    #                  +-> invalid shape or missing key -> raise
    #                  +-> seen key -> skip duplicate
    #                  `-> first key -> write record and count it
    if not isinstance(key_name, str):
        raise TypeError("key_name must be a string")
    if not key_name.strip():
        raise ValueError("key_name cannot be empty")

    seen_keys = set()
    written_count = 0

    with open(input_path, "r", encoding="utf-8") as input_file:
        with open(output_path, "w", encoding="utf-8") as output_file:
            for line_number, raw_line in enumerate(input_file, start=1):
                line = raw_line.strip()
                if not line:
                    continue

                try:
                    record = json.loads(line)
                except json.JSONDecodeError as error:
                    raise ValueError(f"invalid JSON on line {line_number}") from error

                if not isinstance(record, dict):
                    raise ValueError(f"record on line {line_number} must be a JSON object")
                if key_name not in record:
                    raise ValueError(f"missing key '{key_name}' on line {line_number}")
                if not isinstance(record[key_name], str):
                    raise ValueError(f"key '{key_name}' on line {line_number} must be a string")

                key_value = record[key_name]
                if key_value in seen_keys:
                    continue

                seen_keys.add(key_value)
                output_file.write(json.dumps(record, ensure_ascii=False))
                output_file.write("\n")
                written_count += 1

    return written_count


if __name__ == "__main__":
    written = deduplicate_jsonl_records_by_key("events.jsonl", "unique-events.jsonl", "id")
    print(written)