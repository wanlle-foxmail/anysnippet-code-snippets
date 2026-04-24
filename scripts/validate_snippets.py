from __future__ import annotations

import json
import re
import sys
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any


REQUIRED_FIELDS = {
    "schema_version",
    "slug",
    "use_case_key",
    "title",
    "language",
    "summary",
    "tags",
    "highlights",
    "repo_path",
    "entry_file",
    "source_files",
    "test_files",
    "dependencies",
    "test_environment",
    "test_type",
    "test_framework",
    "test_case_count",
    "verification_workdir",
    "verification_command",
    "verification_setup",
    "verification_notes",
    "status",
}

STRING_FIELDS = {
    "slug",
    "use_case_key",
    "title",
    "language",
    "summary",
    "repo_path",
    "entry_file",
    "test_type",
    "test_framework",
    "verification_workdir",
    "verification_command",
    "verification_notes",
    "status",
}

LIST_FIELDS = {
    "tags",
    "highlights",
    "source_files",
    "test_files",
    "dependencies",
    "verification_setup",
}

STATUS_VALUES = {"Draft", "Verified"}
TEST_TYPE_VALUES = {"IntegrationTest", "ManualVerification", "SmokeTest", "UnitTest"}
NON_EMPTY_STRING_FIELDS = STRING_FIELDS
SLUG_PATTERN = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
CJK_PATTERN = re.compile(r"[\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff]")


@dataclass(frozen=True)
class ValidationIssue:
    path: str
    message: str


@dataclass
class ValidationResult:
    errors: list[ValidationIssue] = field(default_factory=list)
    warnings: list[ValidationIssue] = field(default_factory=list)

    def add_error(self, path: Path, message: str) -> None:
        self.errors.append(ValidationIssue(path=path.as_posix(), message=message))

    def add_warning(self, path: Path, message: str) -> None:
        self.warnings.append(ValidationIssue(path=path.as_posix(), message=message))


def validate_repository(root: Path) -> ValidationResult:
    result = ValidationResult()
    snippets_root = root / "snippets"
    if not snippets_root.exists():
        return result

    seen_slugs: dict[str, Path] = {}
    seen_language_use_cases: dict[tuple[str, str], Path] = {}

    for snippet_dir in iter_snippet_dirs(snippets_root):
        metadata = load_metadata(snippet_dir, result)
        validate_required_files(snippet_dir, result)
        validate_readme_content(snippet_dir, result)
        if metadata is None:
            continue

        validate_metadata(snippet_dir, root, metadata, result)

        slug = metadata.get("slug")
        if isinstance(slug, str):
            existing_path = seen_slugs.get(slug)
            if existing_path is None:
                seen_slugs[slug] = snippet_dir
            else:
                result.add_error(snippet_dir, f"Duplicate slug: {slug} (already used by {existing_path.relative_to(root).as_posix()})")

        language = metadata.get("language")
        use_case_key = metadata.get("use_case_key")
        if isinstance(language, str) and isinstance(use_case_key, str):
            dedup_key = (language, use_case_key)
            existing_path = seen_language_use_cases.get(dedup_key)
            if existing_path is None:
                seen_language_use_cases[dedup_key] = snippet_dir
            else:
                result.add_error(
                    snippet_dir,
                    "Duplicate language/use_case_key: "
                    f"{language}/{use_case_key} (already used by {existing_path.relative_to(root).as_posix()})",
                )

    return result


def iter_snippet_dirs(snippets_root: Path) -> list[Path]:
    snippet_dirs: list[Path] = []
    for language_dir in sorted(snippets_root.iterdir()):
        if not language_dir.is_dir():
            continue
        for snippet_dir in sorted(language_dir.iterdir()):
            if snippet_dir.is_dir():
                snippet_dirs.append(snippet_dir)
    return snippet_dirs


def load_metadata(snippet_dir: Path, result: ValidationResult) -> dict[str, Any] | None:
    metadata_path = snippet_dir / "snippet.json"
    if not metadata_path.exists():
        return None

    try:
        raw_data = json.loads(metadata_path.read_text(encoding="utf-8"))
    except (json.JSONDecodeError, UnicodeDecodeError, OSError) as error:
        result.add_error(metadata_path, f"Failed to read or parse JSON: {error}")
        return None

    if not isinstance(raw_data, dict):
        result.add_error(metadata_path, "snippet.json must contain a JSON object")
        return None

    return raw_data


def validate_required_files(snippet_dir: Path, result: ValidationResult) -> None:
    for required_name in ("snippet.json", "README.md"):
        required_path = snippet_dir / required_name
        if not required_path.exists():
            result.add_error(snippet_dir, f"Missing required file: {required_name}")


def validate_metadata(snippet_dir: Path, root: Path, metadata: dict[str, Any], result: ValidationResult) -> None:
    metadata_path = snippet_dir / "snippet.json"

    missing_fields = sorted(REQUIRED_FIELDS - metadata.keys())
    for field_name in missing_fields:
        result.add_error(metadata_path, f"Missing required field: {field_name}")

    validate_field_types(metadata_path, metadata, result)
    validate_schema_rules(snippet_dir, root, metadata, result)
    validate_referenced_files(snippet_dir, metadata_path, metadata, result)
    validate_english_only(metadata_path, metadata, result)
    validate_test_count_rules(metadata_path, metadata, result)


def validate_field_types(metadata_path: Path, metadata: dict[str, Any], result: ValidationResult) -> None:
    for field_name in STRING_FIELDS:
        if field_name in metadata and not isinstance(metadata[field_name], str):
            result.add_error(metadata_path, f"Field {field_name} must be a string")
        if field_name in metadata and isinstance(metadata[field_name], str) and not metadata[field_name].strip():
            result.add_error(metadata_path, f"Field {field_name} cannot be empty")

    for field_name in LIST_FIELDS:
        if field_name in metadata and not isinstance(metadata[field_name], list):
            result.add_error(metadata_path, f"Field {field_name} must be an array")
            continue
        if field_name in metadata:
            for item in metadata[field_name]:
                if not isinstance(item, str):
                    result.add_error(metadata_path, f"Field {field_name} must contain only strings")
                    break
                if field_name in {"source_files", "test_files"} and not item.strip():
                    result.add_error(metadata_path, f"Field {field_name} cannot contain empty strings")
                    break

    if "test_environment" in metadata and not isinstance(metadata["test_environment"], dict):
        result.add_error(metadata_path, "Field test_environment must be an object")
    elif "test_environment" in metadata:
        test_environment = metadata["test_environment"]
        runtime = test_environment.get("runtime")
        supported_os = test_environment.get("os")

        if not isinstance(runtime, str) or not runtime.strip():
            result.add_error(metadata_path, "Field test_environment.runtime must be a non-empty string")
        if not isinstance(supported_os, list) or not supported_os or not all(isinstance(item, str) and item.strip() for item in supported_os):
            result.add_error(metadata_path, "Field test_environment.os must be a non-empty array of strings")

    if "test_case_count" in metadata:
        test_case_count = metadata["test_case_count"]
        if isinstance(test_case_count, bool) or not isinstance(test_case_count, int) or test_case_count < 0:
            result.add_error(metadata_path, "Field test_case_count must be a non-negative integer")

    if "test_count_rationale" in metadata and not isinstance(metadata["test_count_rationale"], str):
        result.add_error(metadata_path, "Field test_count_rationale must be a string")


def validate_schema_rules(snippet_dir: Path, root: Path, metadata: dict[str, Any], result: ValidationResult) -> None:
    metadata_path = snippet_dir / "snippet.json"

    schema_version = metadata.get("schema_version")
    if isinstance(schema_version, bool) or not isinstance(schema_version, int) or schema_version != 1:
        result.add_error(metadata_path, "schema_version must be 1")

    status = metadata.get("status")
    if isinstance(status, str) and status not in STATUS_VALUES:
        result.add_error(metadata_path, f"status must be one of: {', '.join(sorted(STATUS_VALUES))}")

    test_type = metadata.get("test_type")
    if isinstance(test_type, str) and test_type not in TEST_TYPE_VALUES:
        result.add_warning(metadata_path, f"test_type should use a canonical value such as: {', '.join(sorted(TEST_TYPE_VALUES))}")

    slug = metadata.get("slug")
    if isinstance(slug, str):
        if not SLUG_PATTERN.fullmatch(slug):
            result.add_error(metadata_path, f"slug must be lowercase kebab-case: {slug}")
        if slug != snippet_dir.name:
            result.add_error(metadata_path, f"slug does not match folder name: expected {snippet_dir.name}, got {slug}")

    repo_path = metadata.get("repo_path")
    if isinstance(repo_path, str):
        actual_repo_path = snippet_dir.relative_to(root).as_posix()
        if repo_path != actual_repo_path:
            result.add_error(metadata_path, f"repo_path does not match folder path: expected {actual_repo_path}, got {repo_path}")


def validate_referenced_files(
    snippet_dir: Path,
    metadata_path: Path,
    metadata: dict[str, Any],
    result: ValidationResult,
) -> None:
    entry_file = metadata.get("entry_file")
    if isinstance(entry_file, str):
        validate_relative_file_exists(snippet_dir, metadata_path, entry_file, "entry_file", result)

    for field_name in ("source_files", "test_files"):
        files = metadata.get(field_name)
        if not isinstance(files, list):
            continue
        for relative_path in files:
            if isinstance(relative_path, str) and relative_path.strip():
                validate_relative_file_exists(snippet_dir, metadata_path, relative_path, field_name, result)


def validate_relative_file_exists(
    snippet_dir: Path,
    metadata_path: Path,
    relative_path: str,
    field_name: str,
    result: ValidationResult,
) -> None:
    snippet_root = snippet_dir.resolve(strict=False)
    file_path = (snippet_dir / relative_path).resolve(strict=False)
    try:
        file_path.relative_to(snippet_root)
    except ValueError:
        result.add_error(metadata_path, f"{field_name} must not escape snippet directory: {relative_path}")
        return

    if not file_path.exists():
        result.add_error(metadata_path, f"{field_name} references a file that does not exist: {relative_path}")


def validate_english_only(metadata_path: Path, metadata: dict[str, Any], result: ValidationResult) -> None:
    for field_name, value in metadata.items():
        if contains_cjk_text(value):
            result.add_error(metadata_path, f"Field {field_name} must be English-only for public metadata")


def validate_readme_content(snippet_dir: Path, result: ValidationResult) -> None:
    readme_path = snippet_dir / "README.md"
    if not readme_path.exists():
        return

    try:
        readme_text = readme_path.read_text(encoding="utf-8")
    except (UnicodeDecodeError, OSError) as error:
        result.add_error(readme_path, f"Failed to read README.md: {error}")
        return

    if CJK_PATTERN.search(remove_fenced_code_blocks(readme_text)):
        result.add_error(readme_path, "README.md must be English-only for public snippets")


def validate_test_count_rules(metadata_path: Path, metadata: dict[str, Any], result: ValidationResult) -> None:
    status = metadata.get("status")
    test_case_count = metadata.get("test_case_count")
    rationale = metadata.get("test_count_rationale")

    if status != "Verified" or isinstance(test_case_count, bool) or not isinstance(test_case_count, int):
        return

    if test_case_count < 5 and (not isinstance(rationale, str) or not rationale.strip()):
        result.add_error(
            metadata_path,
            "Verified snippets with fewer than 5 test cases must include test_count_rationale",
        )

    if test_case_count > 10:
        result.add_warning(
            metadata_path,
            "test_case_count is above the preferred 5-10 range; confirm the larger suite is justified",
        )


def contains_cjk_text(value: Any) -> bool:
    if isinstance(value, str):
        return bool(CJK_PATTERN.search(value))
    if isinstance(value, list):
        return any(contains_cjk_text(item) for item in value)
    if isinstance(value, dict):
        return any(contains_cjk_text(item) for item in value.values())
    return False


def remove_fenced_code_blocks(content: str) -> str:
    return re.sub(r"```.*?```", "", content, flags=re.DOTALL)


def main(argv: list[str] | None = None) -> int:
    args = argv if argv is not None else sys.argv[1:]
    root = Path(args[0]).resolve() if args else Path(__file__).resolve().parent.parent
    result = validate_repository(root)

    for issue in result.errors:
        print(f"ERROR {issue.path}: {issue.message}")
    for issue in result.warnings:
        print(f"WARNING {issue.path}: {issue.message}")

    if result.errors:
        print(f"Validation failed with {len(result.errors)} error(s) and {len(result.warnings)} warning(s).")
        return 1

    print(f"Validation passed with {len(result.warnings)} warning(s).")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())