from collections.abc import Iterable, Iterator


def batch_records_by_total_byte_size(records: Iterable[str], max_batch_bytes: int) -> Iterator[list[str]]:
    """Yield batches of string records whose total UTF-8 byte size stays under a limit."""
    # Flow:
    #   records -> measure each record in UTF-8 bytes
    #              |
    #              +-> single record too large -> raise
    #              +-> next record would overflow -> yield batch and start a new one
    #              `-> fits current batch -> keep accumulating, then yield the tail batch
    if isinstance(records, str):
        raise TypeError("records must be an iterable of strings")
    if isinstance(max_batch_bytes, bool) or not isinstance(max_batch_bytes, int) or max_batch_bytes <= 0:
        raise ValueError("max_batch_bytes must be a positive integer")

    current_batch = []
    current_batch_bytes = 0

    for record in records:
        if not isinstance(record, str):
            raise TypeError("each record must be a string")

        record_bytes = len(record.encode("utf-8"))
        if record_bytes > max_batch_bytes:
            raise ValueError("single record exceeds max_batch_bytes")

        if current_batch and current_batch_bytes + record_bytes > max_batch_bytes:
            yield current_batch
            current_batch = [record]
            current_batch_bytes = record_bytes
            continue

        current_batch.append(record)
        current_batch_bytes += record_bytes

    if current_batch:
        yield current_batch


if __name__ == "__main__":
    result = list(batch_records_by_total_byte_size(["alpha", "beta", "gamma"], 9))
    print(result)