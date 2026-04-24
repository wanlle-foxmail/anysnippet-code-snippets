def merge_overlapping_time_ranges(ranges: list[tuple[int, int]]) -> list[tuple[int, int]]:
    """Sort integer time ranges by start time and merge overlapping intervals."""
    # Flow:
    #   ranges -> validate and normalize pairs
    #             |
    #             +-> sort by start time
    #                  |
    #                  +-> overlap? yes -> extend last merged range
    #                  |          no -> append a new merged range
    #                  +-> return merged ranges
    if not isinstance(ranges, list):
        raise TypeError("ranges must be a list")

    normalized_ranges = []
    for time_range in ranges:
        if not isinstance(time_range, tuple) or len(time_range) != 2:
            raise TypeError("each range must be a 2-item tuple")

        start, end = time_range
        if isinstance(start, bool) or not isinstance(start, int):
            raise TypeError("range start must be an integer")
        if isinstance(end, bool) or not isinstance(end, int):
            raise TypeError("range end must be an integer")
        if start > end:
            raise ValueError("range start must be less than or equal to range end")

        normalized_ranges.append((start, end))

    if not normalized_ranges:
        return []

    normalized_ranges.sort(key=lambda item: (item[0], item[1]))
    merged_ranges = [normalized_ranges[0]]

    for start, end in normalized_ranges[1:]:
        last_start, last_end = merged_ranges[-1]
        if start <= last_end:
            merged_ranges[-1] = (last_start, max(last_end, end))
        else:
            merged_ranges.append((start, end))

    return merged_ranges


if __name__ == "__main__":
    result = merge_overlapping_time_ranges([(60, 120), (90, 180), (240, 300)])
    print(result)