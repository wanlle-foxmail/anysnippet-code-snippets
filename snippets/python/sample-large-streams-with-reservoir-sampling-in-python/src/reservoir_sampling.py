import random
from collections.abc import Iterable
from typing import Any, Optional


def reservoir_sample(items: Iterable[Any], sample_size: int, seed: Optional[int] = None) -> list[Any]:
    """Sample a fixed number of items from an iterable without buffering the full input."""
    # Flow:
    #   stream -> first sample_size items fill the reservoir
    #             later items -> draw a replacement index
    #                            |
    #                            +-> inside the reservoir -> replace that slot
    #                            `-> outside the reservoir -> keep the current sample
    if isinstance(sample_size, bool) or not isinstance(sample_size, int) or sample_size <= 0:
        raise ValueError("sample_size must be a positive integer")

    random_generator = random.Random(seed)
    reservoir = []

    for index, item in enumerate(items):
        if index < sample_size:
            reservoir.append(item)
            continue

        replacement_index = random_generator.randint(0, index)
        if replacement_index < sample_size:
            reservoir[replacement_index] = item

    return reservoir


if __name__ == "__main__":
    result = reservoir_sample(range(10_000), 5, seed=7)
    print(result)