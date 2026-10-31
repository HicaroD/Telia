"""Exercise the release gate through its CLI and the GitHub HTTP API boundary."""
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import threading
import unittest
from urllib.parse import parse_qs, urlsplit

ROOT = Path(__file__).resolve().parent.parent
COMMIT = "a" * 40
REPOSITORY = "example/Telia"
TARGETS = ("linux-amd64", "darwin-amd64", "darwin-arm64", "windows-amd64")


def run(run_id):
    return {"id": run_id, "head_sha": COMMIT, "head_branch": "master", "event": "push",
            "head_repository": {"full_name": REPOSITORY},
            "status": "completed", "conclusion": "success"}


class ReleaseSourceTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="telia-release-gate-")
        self.addCleanup(self.temp.cleanup)
        self.output = Path(self.temp.name) / "output"
        self.responses = {
            ("ci.yml/runs", 1): {"workflow_runs": [run(10)]},
            ("native.yml/runs", 1): {"workflow_runs": [run(20)]},
            ("20/artifacts", 1): {"artifacts": [
                {"id": index, "name": "sdl-3.2.28-" + target, "expired": False,
                 "size_in_bytes": 100}
                for index, target in enumerate(TARGETS, 1)]},
        }
        self.requests = []
        owner = self

        class Handler(BaseHTTPRequestHandler):
            def do_GET(self):
                request = urlsplit(self.path)
                query = parse_qs(request.query)
                owner.requests.append((request.path, query, self.headers.get("Authorization")))
                key = request.path.removeprefix("/repos/" + REPOSITORY + "/actions/")
                key = key.removeprefix("workflows/").removeprefix("runs/")
                response = owner.responses.get((key, int(query.get("page", [1])[0])))
                self.send_response(200 if response is not None else 404)
                self.send_header("Content-Type", "application/json")
                self.end_headers()
                self.wfile.write(json.dumps(response).encode())

            def log_message(self, *args):
                pass

        self.server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
        self.thread = threading.Thread(target=self.server.serve_forever,
                                       kwargs={"poll_interval": 0.01}, daemon=True)
        self.thread.start()
        self.addCleanup(self.stop_server)

    def stop_server(self):
        self.server.shutdown()
        self.server.server_close()
        self.thread.join()

    def gate(self):
        return subprocess.run([
            sys.executable, str(ROOT / "scripts/release-source.py"),
            "--repository", REPOSITORY, "--commit", COMMIT,
            "--api-url", "http://127.0.0.1:" + str(self.server.server_port),
            "--output", str(self.output),
        ], env=dict(os.environ, GH_TOKEN="test-token"), text=True, capture_output=True)

    def test_selects_verified_native_artifacts_for_the_exact_commit(self):
        result = self.gate()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(self.output.read_text(),
                         "run_id=20\nsdl_version=3.2.28\ncommit=" + COMMIT + "\n")
        self.assertEqual(len(self.requests), 3)
        for path, query, authorization in self.requests:
            self.assertEqual(authorization, "Bearer test-token")
            if path.endswith("/runs"):
                self.assertEqual(query["head_sha"], [COMMIT])
                self.assertEqual(query["branch"], ["master"])
                self.assertEqual(query["event"], ["push"])

    def assert_blocked(self, message):
        result = self.gate()
        self.assertNotEqual(result.returncode, 0, result.stdout)
        self.assertIn(message, result.stderr)
        self.assertFalse(self.output.exists(), "a blocked release must not expose partial outputs")

    def test_requires_both_workflows_to_succeed(self):
        for workflow in ("ci.yml", "native.yml"):
            for status, conclusion in (("in_progress", None), ("completed", "failure"),
                                       ("completed", "cancelled"), ("completed", "skipped")):
                with self.subTest(workflow=workflow, status=status, conclusion=conclusion):
                    record = self.responses[(workflow + "/runs", 1)]["workflow_runs"][0]
                    record.update(status=status, conclusion=conclusion)
                    self.assert_blocked("is not successful")
                    record.update(status="completed", conclusion="success")

    def test_never_reuses_another_commit_branch_event_or_fork(self):
        for change in ({"head_sha": "b" * 40}, {"head_branch": "feature"},
                       {"event": "pull_request"}, {"event": "workflow_dispatch"},
                       {"head_repository": {"full_name": "fork/Telia"}}):
            with self.subTest(change=change):
                record = run(20)
                record.update(change)
                self.responses[("native.yml/runs", 1)]["workflow_runs"] = [record]
                self.assert_blocked("no master push run")

    def test_does_not_fall_back_to_older_success_when_newer_run_failed(self):
        latest = run(21)
        latest["conclusion"] = "failure"
        self.responses[("native.yml/runs", 1)]["workflow_runs"] = [run(20), latest]
        self.assert_blocked("run 21 is not successful")

    def test_missing_expired_empty_or_duplicate_artifacts_block_release(self):
        original = self.responses[("20/artifacts", 1)]["artifacts"]
        for artifacts in (original[:-1],
                          [dict(original[0], expired=True)] + original[1:],
                          [dict(original[0], size_in_bytes=0)] + original[1:],
                          original + [dict(original[0], id=99)]):
            with self.subTest(artifacts=artifacts):
                self.responses[("20/artifacts", 1)]["artifacts"] = artifacts
                self.assert_blocked("release will not rebuild SDL automatically")

    def test_reads_all_pages_of_workflow_runs_and_artifacts(self):
        self.responses[("native.yml/runs", 1)] = {"workflow_runs": [
            dict(run(index), head_sha="b" * 40) for index in range(100)]}
        self.responses[("native.yml/runs", 2)] = {"workflow_runs": [run(20)]}
        artifacts = self.responses[("20/artifacts", 1)]
        self.responses[("20/artifacts", 1)] = {"artifacts": [
            {"name": "unrelated-" + str(index)} for index in range(100)]}
        self.responses[("20/artifacts", 2)] = artifacts
        result = self.gate()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("run_id=20", self.output.read_text())

    def test_api_failure_blocks_release(self):
        del self.responses[("ci.yml/runs", 1)]
        self.assert_blocked("404")


if __name__ == "__main__":
    unittest.main()
