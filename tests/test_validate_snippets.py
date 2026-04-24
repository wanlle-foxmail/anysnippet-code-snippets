import json
import sys
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parent.parent
if str(REPO_ROOT) not in sys.path:
    sys.path.insert(0, str(REPO_ROOT))

from scripts.validate_snippets import validate_repository


def build_metadata(**overrides):
    metadata = {
        "schema_version": 1,
        "slug": "copy-directory-files-with-python",
        "use_case_key": "copy-directory-files",
        "title": "Copy Directory Files with Python",
        "language": "Python",
        "summary": "Copy files from one directory to another while preserving nested structure.",
        "tags": ["filesystem", "automation", "pathlib", "copy"],
        "highlights": [
            "Preserves nested directory structure",
            "Creates missing target directories automatically",
            "Uses pathlib for readable path handling",
        ],
        "repo_path": "snippets/python/copy-directory-files-with-python",
        "entry_file": "src/copy_directory_files.py",
        "source_files": ["src/copy_directory_files.py"],
        "test_files": ["tests/test_copy_directory_files.py"],
        "dependencies": [],
        "test_environment": {"runtime": "Python 3.12", "os": ["macOS", "Linux", "Windows"]},
        "test_type": "UnitTest",
        "test_framework": "pytest",
        "test_case_count": 7,
        "verification_workdir": ".",
        "verification_command": "pytest -q",
        "verification_setup": [],
        "verification_notes": "Executed from the snippet root.",
        "status": "Verified",
    }
    metadata.update(overrides)
    return metadata


def create_snippet(
    root: Path,
    language: str,
    slug: str,
    metadata=None,
    readme_text="README",
    create_referenced_files=True,
):
    snippet_dir = root / "snippets" / language.lower() / slug
    snippet_dir.mkdir(parents=True, exist_ok=True)

    if metadata is not None:
        (snippet_dir / "snippet.json").write_text(json.dumps(metadata, indent=2), encoding="utf-8")
        if create_referenced_files:
            for relative_path in [metadata["entry_file"], *metadata["source_files"], *metadata["test_files"]]:
                target_path = snippet_dir / relative_path
                target_path.parent.mkdir(parents=True, exist_ok=True)
                target_path.write_text("placeholder", encoding="utf-8")

    if readme_text is not None:
        (snippet_dir / "README.md").write_text(readme_text, encoding="utf-8")

    return snippet_dir


class ValidateSnippetsTests(unittest.TestCase):
    def test_valid_snippet_passes(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata())

            result = validate_repository(root)

            self.assertEqual([], result.errors)

    def test_invalid_json_syntax_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            snippet_dir = create_snippet(root, "python", "copy-directory-files-with-python", metadata=None)
            (snippet_dir / "snippet.json").write_text("{invalid json", encoding="utf-8")

            result = validate_repository(root)

            self.assertTrue(any("Failed to read or parse JSON" in error.message for error in result.errors))

    def test_non_object_json_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            snippet_dir = create_snippet(root, "python", "copy-directory-files-with-python", metadata=None)
            (snippet_dir / "snippet.json").write_text("[]", encoding="utf-8")

            result = validate_repository(root)

            self.assertTrue(any("must contain a JSON object" in error.message for error in result.errors))

    def test_missing_required_field_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            metadata = build_metadata()
            del metadata["slug"]
            create_snippet(root, "python", "copy-directory-files-with-python", metadata)

            result = validate_repository(root)

            self.assertTrue(any("Missing required field: slug" in error.message for error in result.errors))

    def test_invalid_schema_version_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(schema_version=2))

            result = validate_repository(root)

            self.assertTrue(any("schema_version must be 1" in error.message for error in result.errors))

    def test_bool_schema_version_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(schema_version=True))

            result = validate_repository(root)

            self.assertTrue(any("schema_version must be 1" in error.message for error in result.errors))

    def test_invalid_status_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(status="Published"))

            result = validate_repository(root)

            self.assertTrue(any("status must be one of" in error.message for error in result.errors))

    def test_duplicate_slug_fails_across_languages(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "calculate-large-file-hashes-with-python", build_metadata(
                slug="calculate-large-file-hashes-with-python",
                use_case_key="calculate-large-file-hashes",
                repo_path="snippets/python/calculate-large-file-hashes-with-python",
            ))
            create_snippet(root, "go", "calculate-large-file-hashes-with-python", build_metadata(
                slug="calculate-large-file-hashes-with-python",
                use_case_key="calculate-large-file-hashes",
                language="Go",
                title="Calculate Large File Hashes with Go",
                repo_path="snippets/go/calculate-large-file-hashes-with-python",
            ))

            result = validate_repository(root)

            self.assertTrue(any("Duplicate slug" in error.message for error in result.errors))

    def test_invalid_slug_format_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "Bad_Slug", build_metadata(slug="Bad_Slug"))

            result = validate_repository(root)

            self.assertTrue(any("slug must be lowercase kebab-case" in error.message for error in result.errors))

    def test_duplicate_language_use_case_key_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata())
            create_snippet(root, "python", "copy-directory-tree-with-python", build_metadata(
                slug="copy-directory-tree-with-python",
                repo_path="snippets/python/copy-directory-tree-with-python",
            ))

            result = validate_repository(root)

            self.assertTrue(any("Duplicate language/use_case_key" in error.message for error in result.errors))

    def test_cross_language_same_use_case_key_with_unique_slugs_passes(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "calculate-large-file-hashes-with-python", build_metadata(
                slug="calculate-large-file-hashes-with-python",
                use_case_key="calculate-large-file-hashes",
                repo_path="snippets/python/calculate-large-file-hashes-with-python",
                title="Calculate Large File Hashes with Python",
            ))
            create_snippet(root, "go", "calculate-large-file-hashes-with-go", build_metadata(
                slug="calculate-large-file-hashes-with-go",
                use_case_key="calculate-large-file-hashes",
                language="Go",
                title="Calculate Large File Hashes with Go",
                repo_path="snippets/go/calculate-large-file-hashes-with-go",
            ))

            result = validate_repository(root)

            self.assertEqual([], result.errors)

    def test_repo_path_mismatch_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(
                repo_path="snippets/python/not-the-real-folder",
            ))

            result = validate_repository(root)

            self.assertTrue(any("repo_path does not match folder path" in error.message for error in result.errors))

    def test_verified_low_test_count_requires_rationale(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "tiny-snippet-with-python", build_metadata(
                slug="tiny-snippet-with-python",
                use_case_key="tiny-snippet",
                title="Tiny Snippet with Python",
                repo_path="snippets/python/tiny-snippet-with-python",
                test_case_count=3,
            ))

            result = validate_repository(root)

            self.assertTrue(any("test_count_rationale" in error.message for error in result.errors))

    def test_bool_test_case_count_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(test_case_count=True))

            result = validate_repository(root)

            self.assertTrue(any("test_case_count must be a non-negative integer" in error.message for error in result.errors))
            self.assertFalse(any("test_count_rationale" in error.message for error in result.errors))

    def test_draft_low_test_count_does_not_require_rationale(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "tiny-draft-snippet-with-python", build_metadata(
                slug="tiny-draft-snippet-with-python",
                use_case_key="tiny-draft-snippet",
                title="Tiny Draft Snippet with Python",
                repo_path="snippets/python/tiny-draft-snippet-with-python",
                test_case_count=1,
                status="Draft",
            ))

            result = validate_repository(root)

            self.assertEqual([], result.errors)

    def test_missing_referenced_file_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(
                root,
                "python",
                "copy-directory-files-with-python",
                build_metadata(),
                create_referenced_files=False,
            )

            result = validate_repository(root)

            self.assertTrue(any("does not exist" in error.message for error in result.errors))

    def test_referenced_files_must_not_escape_snippet_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            metadata = build_metadata(
                entry_file="../outside.py",
                source_files=["src/copy_directory_files.py"],
                test_files=["tests/test_copy_directory_files.py"],
            )
            snippet_dir = create_snippet(
                root,
                "python",
                "copy-directory-files-with-python",
                metadata,
                create_referenced_files=False,
            )
            (snippet_dir / "src" / "copy_directory_files.py").parent.mkdir(parents=True, exist_ok=True)
            (snippet_dir / "src" / "copy_directory_files.py").write_text("placeholder", encoding="utf-8")
            (snippet_dir / "tests" / "test_copy_directory_files.py").parent.mkdir(parents=True, exist_ok=True)
            (snippet_dir / "tests" / "test_copy_directory_files.py").write_text("placeholder", encoding="utf-8")
            (snippet_dir.parent / "outside.py").write_text("outside", encoding="utf-8")

            result = validate_repository(root)

            self.assertTrue(any("must not escape snippet directory" in error.message for error in result.errors))

    def test_source_files_reject_empty_strings(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            metadata = build_metadata(source_files=[""], test_files=["tests/test_copy_directory_files.py"])
            snippet_dir = create_snippet(
                root,
                "python",
                "copy-directory-files-with-python",
                metadata,
                create_referenced_files=False,
            )
            (snippet_dir / "src" / "copy_directory_files.py").parent.mkdir(parents=True, exist_ok=True)
            (snippet_dir / "src" / "copy_directory_files.py").write_text("placeholder", encoding="utf-8")
            (snippet_dir / "tests" / "test_copy_directory_files.py").parent.mkdir(parents=True, exist_ok=True)
            (snippet_dir / "tests" / "test_copy_directory_files.py").write_text("placeholder", encoding="utf-8")

            result = validate_repository(root)

            self.assertTrue(any("source_files cannot contain empty strings" in error.message for error in result.errors))

    def test_verified_low_test_count_with_rationale_passes(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "tiny-snippet-with-python", build_metadata(
                slug="tiny-snippet-with-python",
                use_case_key="tiny-snippet",
                title="Tiny Snippet with Python",
                repo_path="snippets/python/tiny-snippet-with-python",
                test_case_count=3,
                test_count_rationale="The behavior surface only has three meaningful scenarios.",
            ))

            result = validate_repository(root)

            self.assertEqual([], result.errors)

    def test_verified_low_test_count_with_empty_rationale_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "tiny-snippet-with-python", build_metadata(
                slug="tiny-snippet-with-python",
                use_case_key="tiny-snippet",
                title="Tiny Snippet with Python",
                repo_path="snippets/python/tiny-snippet-with-python",
                test_case_count=3,
                test_count_rationale="   ",
            ))

            result = validate_repository(root)

            self.assertTrue(any("test_count_rationale" in error.message for error in result.errors))

    def test_missing_readme_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(
                root,
                "python",
                "copy-directory-files-with-python",
                build_metadata(),
                readme_text=None,
            )

            result = validate_repository(root)

            self.assertTrue(any("Missing required file: README.md" in error.message for error in result.errors))

    def test_missing_snippet_json_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", metadata=None)

            result = validate_repository(root)

            self.assertTrue(any("Missing required file: snippet.json" in error.message for error in result.errors))

    def test_public_metadata_rejects_chinese_text(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(
                summary="复制目录中的文件到另一个目录。",
            ))

            result = validate_repository(root)

            self.assertTrue(any("must be English-only" in error.message for error in result.errors))

    def test_readme_rejects_chinese_text(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(
                root,
                "python",
                "copy-directory-files-with-python",
                build_metadata(),
                readme_text="# Copy Directory Files\n\n这是一个片段。",
            )

            result = validate_repository(root)

            self.assertTrue(any("README.md must be English-only" in error.message for error in result.errors))

    def test_readme_allows_cjk_inside_code_blocks(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(
                root,
                "python",
                "copy-directory-files-with-python",
                build_metadata(),
                readme_text=(
                    "# Copy Directory Files\n\n"
                    "Use this snippet to copy files.\n\n"
                    "```python\n"
                    "print(\"你好\")\n"
                    "```\n"
                ),
            )

            result = validate_repository(root)

            self.assertEqual([], result.errors)

    def test_slug_must_match_folder_name(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(
                slug="different-slug-with-python",
            ))

            result = validate_repository(root)

            self.assertTrue(any("slug does not match folder name" in error.message for error in result.errors))

    def test_empty_required_string_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root = Path(tmp_dir)
            create_snippet(root, "python", "copy-directory-files-with-python", build_metadata(
                verification_command="   ",
            ))

            result = validate_repository(root)

            self.assertTrue(any("cannot be empty" in error.message for error in result.errors))


if __name__ == "__main__":
    unittest.main()