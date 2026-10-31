#!/usr/bin/env python3
"""Verify successful master CI and native artifacts for an exact release commit."""
import argparse
import json
import os
from pathlib import Path
import re
import sys
from urllib.error import HTTPError, URLError
from urllib.parse import urlencode
from urllib.request import Request, urlopen

ROOT = Path(__file__).resolve().parent.parent


def pages(api_url, token, path, key, **query):
    for page in range(1, 101):
        url = api_url.rstrip("/") + path + "?" + urlencode(dict(query, per_page=100, page=page))
        request = Request(url, headers={
            "Authorization": "Bearer " + token,
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
            "User-Agent": "telia-release-gate",
        })
        with urlopen(request, timeout=30) as response:
            items = json.load(response)[key]
        yield from items
        if len(items) < 100:
            return
    raise ValueError("GitHub API pagination limit reached")


def verified_run(api_url, token, repository, commit, workflow):
    runs = pages(api_url, token,
                 "/repos/{}/actions/workflows/{}/runs".format(repository, workflow),
                 "workflow_runs", branch="master", event="push", head_sha=commit)
    matching = [run for run in runs
                if run.get("head_sha") == commit and run.get("head_branch") == "master"
                and run.get("event") == "push"
                and run.get("head_repository", {}).get("full_name", "").lower() == repository.lower()]
    if not matching:
        raise ValueError("{} has no master push run for {}. Push the commit to master and wait for CI."
                         .format(workflow, commit))
    latest = max(matching, key=lambda run: int(run["id"]))
    if latest.get("status") != "completed" or latest.get("conclusion") != "success":
        raise ValueError("{} run {} is not successful ({} / {}). Wait for it or rerun the failed checks."
                         .format(workflow, latest["id"], latest.get("status"), latest.get("conclusion")))
    return int(latest["id"])


def verify_artifacts(api_url, token, repository, run_id, manifest):
    artifacts = list(pages(api_url, token,
                          "/repos/{}/actions/runs/{}/artifacts".format(repository, run_id), "artifacts"))
    for target in manifest["targets"]:
        name = "sdl-{}-{}".format(manifest["version"], target)
        matches = [artifact for artifact in artifacts if artifact.get("name") == name]
        if (len(matches) != 1 or matches[0].get("expired") is not False
                or matches[0].get("size_in_bytes", 0) <= 0):
            raise ValueError("Native run {} has no unique, available artifact {}. "
                             "Rerun Native vendor packages for the same master commit; "
                             "release will not rebuild SDL automatically.".format(run_id, name))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repository", required=True)
    parser.add_argument("--commit", required=True)
    parser.add_argument("--api-url", default=os.environ.get("GITHUB_API_URL", "https://api.github.com"))
    parser.add_argument("--output", type=Path, required=True, help="GitHub Actions output file")
    args = parser.parse_args()
    token = os.environ.get("GH_TOKEN") or os.environ.get("GITHUB_TOKEN")
    try:
        if not token:
            raise ValueError("GH_TOKEN or GITHUB_TOKEN is required to read workflow runs and artifacts")
        if not re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", args.repository):
            raise ValueError("invalid repository; expected owner/name")
        if not re.fullmatch(r"(?:[0-9a-f]{40}|[0-9a-f]{64})", args.commit):
            raise ValueError("release commit must be a full Git commit SHA")
        manifest = json.loads((ROOT / "base/vendor/sdl/native.json").read_text())
        if not re.fullmatch(r"[A-Za-z0-9.+-]+", manifest["version"]):
            raise ValueError("invalid SDL version in native manifest")
        verified_run(args.api_url, token, args.repository, args.commit, "ci.yml")
        native_run = verified_run(args.api_url, token, args.repository, args.commit, "native.yml")
        verify_artifacts(args.api_url, token, args.repository, native_run, manifest)
        with args.output.open("a") as output:
            output.write("run_id={}\nsdl_version={}\ncommit={}\n"
                         .format(native_run, manifest["version"], args.commit))
    except (ValueError, KeyError, TypeError, OSError, HTTPError, URLError) as error:
        print("Release blocked: " + str(error), file=sys.stderr)
        return 1
    print("Verified master commit {} and native artifacts from run {}.".format(args.commit, native_run))
    return 0


if __name__ == "__main__":
    sys.exit(main())
