"""Offline installer checks: downloads are replaced with local release fixtures."""
import hashlib
import io
import os
from pathlib import Path
import subprocess
import tarfile
import tempfile
import unittest

ROOT = Path(__file__).resolve().parent.parent


@unittest.skipIf(os.name == "nt", "POSIX installer; Windows installer is checked on Windows CI")
class InstallerTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="telia-installer-test-")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name).resolve()
        self.tools = self.root / "tools"
        self.tools.mkdir()
        self.fixtures = self.root / "fixtures"
        self.fixtures.mkdir()
        self.env = dict(os.environ, PATH=str(self.tools) + os.pathsep + os.environ["PATH"],
                        HOME=str(self.root), TELIA_VERSION="v0.1.0",
                        TELIA_INSTALL_ROOT=str(self.root / "install"),
                        TELIA_BIN_DIR=str(self.root / "bin"), FIXTURES=str(self.fixtures))
        self.tool("uname", '#!/bin/sh\ncase "$1" in -s) echo Darwin;; -m) echo arm64;; esac\n')
        self.tool("curl", '''#!/bin/sh
case "$*" in
    *releases/latest*) printf 'https://github.com/HicaroD/Telia/releases/tag/v0.1.0'; exit 0 ;;
esac
while [ "$#" -gt 0 ]; do
    case "$1" in
        https://*) asset=${1##*/} ;;
        -o) shift; output=$1 ;;
    esac
    shift
done
cp "$FIXTURES/$asset" "$output"
''')
        self.archive = self.fixtures / "telia-darwin-arm64.tar.gz"
        with tarfile.open(self.archive, "w:gz") as archive:
            files = [("telia", b"#!/bin/sh\necho telia-test\n", 0o755),
                     ("base/std/io/io.t", b"fixture", 0o644),
                     ("base/vendor/example/native.json", b"fixture", 0o644)]
            inventory = "".join(hashlib.sha256(data).hexdigest() + "  " + name + "\n"
                                for name, data, _ in files).encode()
            for name, data, mode in files + [("BUNDLE_SHA256SUMS", inventory, 0o644)]:
                entry = tarfile.TarInfo(name)
                entry.size, entry.mode = len(data), mode
                archive.addfile(entry, io.BytesIO(data))
        digest = hashlib.sha256(self.archive.read_bytes()).hexdigest()
        (self.fixtures / "SHA256SUMS").write_text(digest + "  " + self.archive.name + "\n")

    def tool(self, name, content):
        path = self.tools / name
        path.write_text(content)
        path.chmod(0o755)

    def run_installer(self):
        return subprocess.run(["sh", str(ROOT / "scripts/install.sh")], env=self.env,
                              text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    def test_verified_bundle_and_symlink(self):
        result = self.run_installer()
        self.assertEqual(result.returncode, 0, result.stderr)
        binary = self.root / "bin/telia"
        self.assertTrue(binary.is_symlink())
        self.assertEqual(subprocess.check_output([str(binary)], text=True).strip(), "telia-test")
        self.assertTrue((binary.resolve().parent / "base/std/io/io.t").is_file())
        self.assertNotIn("THE TELIA", result.stdout)  # No banner in redirected output.
        self.assertEqual(list((self.root / "install").glob(".install.*")), [])

    def test_latest_resolution(self):
        self.env["TELIA_VERSION"] = "latest"
        result = self.run_installer()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("v0.1.0", result.stdout)

    def test_checksum_mismatch_installs_nothing(self):
        with self.archive.open("ab") as archive:
            archive.write(b"tampered")
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Checksum mismatch", result.stderr)
        self.assertFalse((self.root / "bin/telia").exists())
        self.assertFalse((self.root / "install/v0.1.0").exists())
        self.assertEqual(list((self.root / "install").iterdir()), [])

    def test_unrelated_binary_is_preserved(self):
        binary = self.root / "bin/telia"
        binary.parent.mkdir()
        binary.write_text("user binary")
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(binary.read_text(), "user binary")

    def test_missing_vendor_file_is_rejected(self):
        self.rewrite_bundle("base/vendor/example/native.json", None)
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Incomplete release bundle", result.stderr)
        self.assertFalse((self.root / "install/v0.1.0").exists())

    def test_corrupt_std_file_is_rejected(self):
        self.rewrite_bundle("base/std/io/io.t", b"corrupt")
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Bundle checksum mismatch", result.stderr)
        self.assertFalse((self.root / "install/v0.1.0").exists())

    def rewrite_bundle(self, name, replacement):
        with tarfile.open(self.archive) as archive:
            files = [(member, archive.extractfile(member).read()) for member in archive.getmembers()]
        with tarfile.open(self.archive, "w:gz") as archive:
            for member, data in files:
                if member.name == name:
                    if replacement is None:
                        continue
                    data = replacement
                    member.size = len(data)
                archive.addfile(member, io.BytesIO(data))
        digest = hashlib.sha256(self.archive.read_bytes()).hexdigest()
        (self.fixtures / "SHA256SUMS").write_text(digest + "  " + self.archive.name + "\n")

    def test_unsupported_platform_fails_before_download(self):
        self.tool("uname", '#!/bin/sh\ncase "$1" in -s) echo Linux;; -m) echo aarch64;; esac\n')
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("No release is available for linux-arm64", result.stderr)

    def test_upgrade_preserves_previous_version(self):
        self.assertEqual(self.run_installer().returncode, 0)
        self.env["TELIA_VERSION"] = "v0.2.0"
        result = self.run_installer()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue((self.root / "install/v0.1.0/telia").exists())
        self.assertEqual((self.root / "bin/telia").resolve(), self.root / "install/v0.2.0/telia")
        self.assertNotEqual(self.run_installer().returncode, 0)


if __name__ == "__main__":
    unittest.main()
