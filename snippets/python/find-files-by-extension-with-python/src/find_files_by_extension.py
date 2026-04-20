import os


def find_files_by_extension(source: str,target: list[str]) -> dict[str, object]:
    """Walk a directory tree and collect files that match the given extensions.

    Args:
        source: Directory path to scan.
        target: Extensions to match, such as [".csv", ".txt"].

    Returns:
        A dictionary with the scanned file count, matched file count, and the
        matched file paths.
    """
    # Scan every file under the source directory and collect the ones that match.
    data = {"count": 0, "hit": 0, "items": []}

    for root, dirs, files in os.walk(source):
        for file in files:
            data["count"] += 1

            # Check whether the current file extension is in the target list.
            _, ext = os.path.splitext(file)
            if ext in target:
                data["hit"] += 1
                filename = os.path.join(root, file)
                data["items"].append(filename)

    return data


if __name__ == "__main__":
    result = find_files_by_extension("Your Directory", [".csv", ".txt"])
    print(result)