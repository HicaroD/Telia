"""Release package discovery and completeness checks without a native toolchain."""
import hashlib
import importlib.util
from pathlib import Path
import tempfile
import unittest

spec = importlib.util.spec_from_file_location("package_release", Path(__file__).with_name("package-release.py"))
release = importlib.util.module_from_spec(spec)
spec.loader.exec_module(release)


class PackageReleaseTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.source = self.root / "vendor"
        self.package = self.source / "second-package"
        self.package.mkdir(parents=True)
        (self.package / "native.json").write_text('''{
          "version": "1", "targets": {
            "linux-amd64": {"sources": ["adapter.c"], "include_dirs": ["."],
              "libraries": ["artifacts/linux-amd64/lib.so"],
              "runtime_files": ["artifacts/linux-amd64/lib.so"]}
          }}''')
        (self.package / "adapter.c").write_text("adapter")
        for target in ("linux-amd64", "darwin-arm64"):
            directory = self.package / "artifacts" / target
            directory.mkdir(parents=True)
            (directory / "lib.so").write_text("library")
        self.destination = self.root / "bundle"

    def test_discovers_packages_and_prunes_other_targets(self):
        release.copy_libraries(self.source, self.destination, "linux-amd64")
        package = self.destination / self.package.name
        self.assertTrue((package / "adapter.c").is_file())
        self.assertTrue((package / "artifacts/linux-amd64/lib.so").is_file())
        self.assertFalse((package / "artifacts/darwin-arm64").exists())

    def test_missing_native_input_fails_packaging(self):
        (self.package / "adapter.c").unlink()
        with self.assertRaisesRegex(ValueError, "missing native artifact"):
            release.copy_libraries(self.source, self.destination, "linux-amd64")

    def test_unsupported_package_target_fails_packaging(self):
        with self.assertRaisesRegex(ValueError, "unsupported"):
            release.copy_libraries(self.source, self.destination, "windows-amd64")

    def test_inventory_covers_every_shipped_file(self):
        release.copy_libraries(self.source, self.destination, "linux-amd64")
        release.write_inventory(self.destination)
        entries = (self.destination / "BUNDLE_SHA256SUMS").read_text().splitlines()
        self.assertEqual(len(entries), 3)
        for entry in entries:
            digest, name = entry.split("  ", 1)
            self.assertEqual(digest, hashlib.sha256((self.destination / name).read_bytes()).hexdigest())
