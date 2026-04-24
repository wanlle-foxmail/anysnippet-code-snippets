from typing import Dict, Iterable, Mapping, Set


def pick_allowed_patch_fields(
    patch_body: Mapping[str, object], allowed_fields: Iterable[str]
) -> Dict[str, object]:
    """Return only allowed top-level fields from a PATCH-style mapping."""
    # Flow:
    #   read the incoming patch body and allowed field list
    #      |
    #      +-> matching top-level keys -> keep them with original values
    #      `-> unknown keys -> ignore them in the returned patch
    if not isinstance(patch_body, Mapping):
        raise TypeError("patch_body must be a mapping")

    allowed_field_names = _normalize_allowed_fields(allowed_fields)
    return {
        key: value
        for key, value in patch_body.items()
        if isinstance(key, str) and key in allowed_field_names
    }


def _normalize_allowed_fields(allowed_fields: Iterable[str]) -> Set[str]:
    normalized = set()
    for field_name in allowed_fields:
        if not isinstance(field_name, str) or field_name.strip() == "":
            raise ValueError("allowed field names must be non-empty strings")
        normalized.add(field_name.strip())
    return normalized


if __name__ == "__main__":
    filtered_patch = pick_allowed_patch_fields(
        {
            "display_name": "Ada",
            "bio": None,
            "role": "admin",
        },
        ["display_name", "bio"],
    )
    print(filtered_patch)