import importlib.util
import tempfile
import unittest
from pathlib import Path


SPEC = importlib.util.spec_from_file_location("check_docs", Path(__file__).with_name("check-docs.py"))
DOCS = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(DOCS)


class DocumentationLinksTest(unittest.TestCase):
    def test_package_index_covers_nested_source_packages(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            source = root / "base" / "std" / "unicode" / "utf8"
            source.mkdir(parents=True)
            (source / "utf8.t").write_text("package utf8\n", encoding="utf-8")
            index = root / "docs" / "stdlib" / "README.md"
            index.parent.mkdir(parents=True)
            index.write_text("# Standard library\n", encoding="utf-8")
            self.assertEqual(len(DOCS.check_packages(root)), 1)
            index.write_text("| `std::unicode::utf8` | UTF-8 |\n", encoding="utf-8")
            self.assertEqual(DOCS.check_packages(root), [])

    def test_anchors_include_duplicates_and_ignore_code(self):
        self.assertEqual(
            DOCS.anchors("# `std::io`\n## Same\n## Same\n```text\n# Fake\n```\n"),
            {"stdio", "same", "same-1"},
        )

    def test_local_files_anchors_and_external_links(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            target = root / "a page.md"
            target.write_text("# Real heading\n", encoding="utf-8")
            page = root / "README.md"
            page.write_text(
                "[ok](a%20page.md#real-heading)\n"
                "[web](https://example.com/missing#anything)\n"
                "```rust\n[fake](not-a-file.md)\n```\n"
                "[bad](missing.md)\n[anchor](a%20page.md#missing)\n",
                encoding="utf-8",
            )
            failures = DOCS.check([page], root)
            self.assertEqual(len(failures), 2)
            self.assertIn("missing target", failures[0])
            self.assertIn("missing heading anchor", failures[1])


if __name__ == "__main__":
    unittest.main()
