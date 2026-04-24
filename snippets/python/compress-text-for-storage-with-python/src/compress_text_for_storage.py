import base64
import gzip


def compress_text(text: str) -> str:
    """Compress text to a base64-encoded gzip string."""
    if not isinstance(text, str):
        raise TypeError("text must be a string")

    text_bytes = text.encode("utf-8")
    compressed_bytes = gzip.compress(text_bytes)
    return base64.b64encode(compressed_bytes).decode("utf-8")


def decompress_text(compressed_text: str) -> str:
    """Decompress a base64-encoded gzip string back to UTF-8 text."""
    if not isinstance(compressed_text, str):
        raise TypeError("compressed_text must be a string")

    try:
        compressed_bytes = base64.b64decode(compressed_text, validate=True)
        text_bytes = gzip.decompress(compressed_bytes)
        return text_bytes.decode("utf-8")
    except (base64.binascii.Error, OSError, UnicodeDecodeError) as error:
        raise ValueError("compressed_text must be a valid base64-encoded gzip string") from error


def main() -> None:
    example_text = '{"page": "/docs", "events": ["view", "view", "view"]}'
    print(compress_text(example_text))


if __name__ == "__main__":
    main()