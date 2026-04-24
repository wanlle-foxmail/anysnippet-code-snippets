import heapq
from collections.abc import Iterator
from contextlib import ExitStack


def merge_sorted_text_files(paths: list[str]) -> Iterator[str]:
    """Yield merged lines from sorted UTF-8 text files without loading them fully."""
    # Flow:
    #   paths -> open files and seed a heap with each first line
    #            |
    #            +-> pop the smallest available line
    #                 |
    #                 +-> read the next line from that same file
    #                      |
    #                      +-> push it back into the heap until all files are exhausted
    if not isinstance(paths, list):
        raise TypeError("paths must be a list")

    with ExitStack() as stack:
        heap = []

        for file_index, path in enumerate(paths):
            if not isinstance(path, str):
                raise TypeError("each path must be a string")

            file_handle = stack.enter_context(open(path, "r", encoding="utf-8"))
            first_line = file_handle.readline()
            if first_line:
                heapq.heappush(heap, (first_line.rstrip("\n"), file_index, file_handle))

        while heap:
            line, file_index, file_handle = heapq.heappop(heap)
            yield line

            next_line = file_handle.readline()
            if next_line:
                heapq.heappush(heap, (next_line.rstrip("\n"), file_index, file_handle))


if __name__ == "__main__":
    for item in merge_sorted_text_files(["part-1.txt", "part-2.txt"]):
        print(item)