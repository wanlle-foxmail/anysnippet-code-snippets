from dataclasses import dataclass
from typing import Iterable, List, Mapping, Set


@dataclass(frozen=True)
class SortField:
    field_name: str
    descending: bool


def parse_sort_query_params(
    query_params: Mapping[str, str], allowed_fields: Iterable[str]
) -> List[SortField]:
    """Parse a comma-separated sort query string against an allowlist."""
    # Flow:
    #   read the sort query value and allowed field list
    #      |
    #      +-> valid allowlisted terms -> return normalized sort fields
    #      `-> unknown or empty terms -> raise ValueError
    raw_sort = query_params.get("sort")
    if raw_sort is None:
        return []
    if not isinstance(raw_sort, str):
        raise ValueError("sort must be a string")

    allowlist = _normalize_allowed_fields(allowed_fields)
    sort_fields = []
    seen_fields = set()

    for raw_term in raw_sort.split(","):
        normalized_term = raw_term.strip()
        if normalized_term == "":
            raise ValueError("sort terms must not be empty")

        descending = normalized_term.startswith("-")
        field_name = normalized_term[1:] if descending else normalized_term
        if field_name not in allowlist:
            raise ValueError(f"sort field is not allowed: {field_name}")
        if field_name in seen_fields:
            raise ValueError(f"sort field must not be repeated: {field_name}")

        seen_fields.add(field_name)
        sort_fields.append(SortField(field_name=field_name, descending=descending))

    return sort_fields


def _normalize_allowed_fields(allowed_fields: Iterable[str]) -> Set[str]:
    normalized = set()
    for field_name in allowed_fields:
        if not isinstance(field_name, str) or field_name.strip() == "":
            raise ValueError("allowed field names must be non-empty strings")
        normalized.add(field_name.strip())
    return normalized


if __name__ == "__main__":
    sort_fields = parse_sort_query_params(
        {"sort": "name,-created_at"},
        ["name", "created_at"],
    )
    print(sort_fields)