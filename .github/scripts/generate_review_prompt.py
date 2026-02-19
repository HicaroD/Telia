import os
import requests

GITHUB_TOKEN = os.environ["GITHUB_TOKEN"]
REPO = os.environ["REPO"]
PR_NUMBER = os.environ["PR_NUMBER"]
REVIEW_ID = os.environ.get("REVIEW_ID", "").strip()

HEADERS = {
    "Authorization": f"Bearer {GITHUB_TOKEN}",
    "Accept": "application/vnd.github+json",
    "X-GitHub-Api-Version": "2022-11-28",
}
BASE = "https://api.github.com"


def fetch_pr():
    r = requests.get(f"{BASE}/repos/{REPO}/pulls/{PR_NUMBER}", headers=HEADERS)
    r.raise_for_status()
    return r.json()


def fetch_reviews():
    r = requests.get(
        f"{BASE}/repos/{REPO}/pulls/{PR_NUMBER}/reviews",
        headers=HEADERS,
        params={"per_page": 100},
    )
    r.raise_for_status()
    return r.json()


def fetch_review_comments(review_id: str):
    """Fetch inline comments scoped to a specific review."""
    r = requests.get(
        f"{BASE}/repos/{REPO}/pulls/{PR_NUMBER}/reviews/{review_id}/comments",
        headers=HEADERS,
        params={"per_page": 100},
    )
    r.raise_for_status()
    return r.json()


def resolve_review(reviews: list) -> dict:
    """Return the target review: by ID if provided, otherwise most recent CHANGES_REQUESTED."""
    if REVIEW_ID:
        match = next((r for r in reviews if str(r["id"]) == REVIEW_ID), None)
        if not match:
            raise ValueError(f"Review ID {REVIEW_ID} not found on PR #{PR_NUMBER}")
        return match

    changes_requested = [r for r in reviews if r["state"] == "CHANGES_REQUESTED"]
    if not changes_requested:
        raise ValueError("No CHANGES_REQUESTED reviews found on this PR")

    return changes_requested[-1]  # most recent


def build_markdown(pr, review, comments) -> str:
    lines = []

    lines.append(f"# PR #{PR_NUMBER} — Agent Change Instructions")
    lines.append(f"branch: `{pr['head']['ref']}` -> `{pr['base']['ref']}`")
    lines.append(f"review: {review['id']} by @{review['user']['login']}")
    lines.append("")
    lines.append("Apply each task below in order. Each task specifies the exact file and line to change.")
    lines.append("")

    if review.get("body", "").strip():
        lines.append("## General feedback")
        lines.append(review["body"].strip())
        lines.append("")

    if not comments:
        lines.append("No inline comments found.")
        return "\n".join(lines)

    by_file: dict[str, list] = {}
    for c in comments:
        by_file.setdefault(c["path"], []).append(c)

    lines.append("## Tasks")
    lines.append("")

    task_num = 1
    for file_path, file_comments in by_file.items():
        for c in file_comments:
            line = c.get("line") or c.get("original_line", "?")
            diff_hunk = c.get("diff_hunk", "").strip()

            lines.append(f"### Task {task_num}")
            lines.append(f"file: `{file_path}` line: {line}")
            lines.append("")

            if diff_hunk:
                lines.append("```diff")
                lines.append(diff_hunk)
                lines.append("```")
                lines.append("")

            lines.append(c["body"].strip())
            lines.append("")

            task_num += 1

    lines.append(f"total tasks: {task_num - 1}")
    return "\n".join(lines)


def post_comment(body: str):
    r = requests.post(
        f"{BASE}/repos/{REPO}/issues/{PR_NUMBER}/comments",
        headers=HEADERS,
        json={"body": body},
    )
    r.raise_for_status()
    print(f"Posted: {r.json()['html_url']}")


if __name__ == "__main__":
    print(f"Fetching PR #{PR_NUMBER} in {REPO} (review_id={REVIEW_ID or 'latest'})...")
    pr = fetch_pr()
    reviews = fetch_reviews()

    review = resolve_review(reviews)
    print(f"Targeting review {review['id']} by @{review['user']['login']}")

    comments = fetch_review_comments(review["id"])
    print(f"{len(comments)} inline comments")

    markdown = build_markdown(pr, review, comments)

    wrapped = f"""<details>
<summary>Agent prompt — review {review['id']} — {len(comments)} task(s)</summary>

{markdown}

</details>"""

    post_comment(wrapped)
