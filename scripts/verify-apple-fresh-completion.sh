#!/usr/bin/env bash
# Read-only completion gate for the single fixed-parent ARM acceptance exception.
# Call only after the independent offline custody/stage/final-byte audits.
set -euo pipefail
[[ $# == 1 ]] || { echo 'usage: verify-apple-fresh-completion.sh EVIDENCE_DIR' >&2; exit 1; }
exec python3 - "$1" <<'PY'
import datetime
import hashlib
import io
import json
import os
from pathlib import Path
import re
import stat
import subprocess
import sys
import tempfile
import time
import zipfile
import zlib

REPO = "Dunky13/takedown"
SOURCE = "ef82381d509aac8b5b02da0f784a8176b1c5b9c0"
PARENT = "27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe"
API = "repos/" + REPO
URL = "https://api.github.com/" + API
PROOF_LIMIT = 1024 * 1024
DEADLINE = time.monotonic() + 180
CRITICAL = (
    "Qualify pinned automation before fresh acceptance",
    "Set up owned Apple ticket transport preference",
    "Install and assess exact final bytes on fresh native runner",
    "Verify owned cleanup before fresh receipt promotion",
    "Promote fresh receipt after verified transport cleanup",
)


class Hold(Exception):
    pass


def require(condition, reason):
    if not condition:
        raise Hold(reason)


def positive(value):
    return type(value) is int and value > 0


def digest(value, length=64):
    return isinstance(value, str) and re.fullmatch("[0-9a-f]{%d}" % length, value)


def pairs(items):
    result = {}
    for key, value in items:
        require(key not in result, "duplicate JSON key")
        result[key] = value
    return result


def parse(raw):
    def invalid_constant(_value):
        raise Hold("invalid JSON constant")
    return json.loads(raw, object_pairs_hook=pairs, parse_constant=invalid_constant)


def api(path, binary=None):
    # Explicit host, GET and finite timeouts; never log gh stderr (SAS URLs).
    args = ["gh", "api", "--hostname", "github.com", "--method", "GET",
            "-H", "X-GitHub-Api-Version: 2022-11-28"]
    if binary:
        args += ["--allow-escape-sequences"]
    if binary == "release":
        args += ["-H", "Accept: application/octet-stream"]
    else:
        args += ["-H", "Accept: application/vnd.github+json"]
    args.append(path)
    environment = os.environ.copy()
    for key in ("GH_DEBUG", "GH_FORCE_TTY"):
        environment.pop(key, None)
    environment["GH_PROMPT_DISABLED"] = "1"
    remaining = DEADLINE - time.monotonic()
    require(remaining > 0, "completion API time budget exhausted")
    with tempfile.TemporaryFile() as output:
        response = subprocess.run(args, stdout=output, stderr=subprocess.DEVNULL,
                                  env=environment, timeout=min(20, remaining), check=False)
        require(response.returncode == 0, "authoritative API read failed; completion held")
        limit = 2 * PROOF_LIMIT if binary else 8 * PROOF_LIMIT
        require(output.tell() <= limit, "API response exceeds completion evidence limit")
        output.seek(0)
        raw = output.read()
    return raw if binary else parse(raw)


def listing(path, key):
    result = []
    ids = set()
    total = None
    # Explicit pagination, consistent totals, no deduplication hiding ambiguity.
    for page in range(1, 21):
        data = api(path + ("&" if "?" in path else "?") + f"per_page=100&page={page}")
        count = data.get("total_count")
        require(type(count) is int and 0 <= count <= 2000, "invalid inventory count")
        if total is None:
            total = count
        require(total == count, "inventory changed during pagination")
        entries = data.get(key)
        require(isinstance(entries, list) and len(entries) <= 100, "invalid inventory page")
        for entry in entries:
            identifier = entry.get("id")
            require(positive(identifier) and identifier not in ids, "duplicate or invalid inventory ID")
            ids.add(identifier)
        result.extend(entries)
        require(len(result) <= total, "inventory count mismatch")
        if len(result) == total:
            return result
        require(len(entries) == 100, "incomplete inventory pagination")
    raise Hold("inventory pagination limit exceeded")


def timestamp(value):
    require(isinstance(value, str), "missing completion timestamp")
    parsed = datetime.datetime.fromisoformat(value.replace("Z", "+00:00"))
    require(parsed.tzinfo is not None, "unqualified completion timestamp")
    return parsed


def completed(item):
    return item.get("status") == "completed" and item.get("conclusion") == "success"


def workflow_path(value, filename, sha):
    base = ".github/workflows/" + filename
    return value in (base, base + "@main", base + "@refs/heads/main", base + "@" + sha)


def verify_job(run_id, attempt, sha):
    run_path = f"{API}/actions/runs/{run_id}"
    run = api(f"{run_path}/attempts/{attempt}")
    require(run.get("id") == run_id and run.get("run_attempt") == attempt,
            "original verification run/attempt mismatch")
    if run.get("status") == "completed":
        require(run.get("conclusion") in ("success", "failure", "cancelled", "neutral", "skipped",
                                          "timed_out", "action_required", "stale", "startup_failure"),
                "unknown original workflow completion state")
    else:
        require(run.get("status") in ("queued", "in_progress", "waiting", "requested", "pending") and
                run.get("conclusion") is None, "unknown original workflow state")
    require(run.get("head_sha") == sha and run.get("head_branch") == "main",
            "original verification automation SHA/branch mismatch")
    require(run.get("event") in ("workflow_dispatch", "workflow_run", "schedule") and
            run.get("pull_requests") == [], "untrusted verification event")
    repository = run.get("repository", {})
    head_repository = run.get("head_repository", {})
    require(repository.get("full_name") == REPO and head_repository.get("full_name") == REPO and
            positive(repository.get("id")) and repository.get("id") == head_repository.get("id") and
            repository.get("fork") is False and head_repository.get("fork") is False,
            "original verification repository mismatch")
    require(run.get("url") == f"{URL}/actions/runs/{run_id}" and
            workflow_path(run.get("path"), "release.yml", sha), "untrusted release workflow path")
    workflow_id = run.get("workflow_id")
    require(positive(workflow_id), "invalid release workflow ID")
    workflow = api(f"{API}/actions/workflows/{workflow_id}")
    require(workflow.get("id") == workflow_id and workflow.get("path") == ".github/workflows/release.yml",
            "release workflow identity mismatch")
    references = run.get("referenced_workflows")
    require(isinstance(references, list), "missing reusable workflow provenance")
    matches = [item for item in references if isinstance(item.get("path"), str) and
               item["path"].startswith(REPO + "/.github/workflows/release-macos.yml@")]
    require(len(matches) == 1 and matches[0].get("sha") == sha and
            workflow_path(matches[0]["path"][len(REPO) + 1:], "release-macos.yml", sha) and
            matches[0].get("ref") == "refs/heads/main",
            "fresh reusable workflow provenance mismatch")
    branch = api(f"{API}/branches/main")
    main_sha = branch.get("commit", {}).get("sha")
    require(branch.get("name") == "main" and branch.get("protected") is True and digest(main_sha, 40),
            "release main is not protected")
    comparison = api(f"{API}/compare/{sha}...{main_sha}?per_page=1")
    require(comparison.get("status") in ("identical", "ahead") and
            comparison.get("merge_base_commit", {}).get("sha") == sha and
            comparison.get("base_commit", {}).get("sha") == sha,
            "verification automation is outside protected main history")
    jobs = listing(f"{run_path}/attempts/{attempt}/jobs", "jobs")
    selected = [job for job in jobs if job.get("name") == "macos / acceptance (macos-15, arm64)"]
    require(len(selected) == 1, "original fresh ARM job missing or ambiguous")
    job = selected[0]
    require(job.get("run_id") == run_id and job.get("run_attempt") == attempt and
            job.get("head_sha") == sha and job.get("head_branch") == "main" and
            job.get("run_url") == f"{URL}/actions/runs/{run_id}" and
            job.get("url") == f"{URL}/actions/jobs/{job['id']}" and
            job.get("labels") == ["macos-15"], "original fresh ARM job binding mismatch")
    require(completed(job), "original fresh ARM job is not completed/success; spent attempt held")
    start, end = timestamp(job.get("started_at")), timestamp(job.get("completed_at"))
    require(start <= end, "invalid original job time range")
    steps = job.get("steps")
    require(isinstance(steps, list), "missing original job steps")
    numbers = [step.get("number") for step in steps]
    require(all(positive(number) for number in numbers) and len(numbers) == len(set(numbers)),
            "ambiguous job step numbers")
    previous = 0
    for name in CRITICAL:
        selected_steps = [step for step in steps if step.get("name") == name]
        require(len(selected_steps) == 1 and completed(selected_steps[0]),
                "critical fresh qualification/cleanup/promotion step did not succeed")
        number = selected_steps[0]["number"]
        require(number > previous, "critical fresh job steps out of order")
        previous = number
    return repository["id"], start, end


def verify_archive(artifact, run_id, sha, repository_id, start, end, proof):
    identifier = artifact["id"]
    require(artifact.get("url") == f"{URL}/actions/artifacts/{identifier}" and
            artifact.get("archive_download_url") == f"{URL}/actions/artifacts/{identifier}/zip",
            "foreign Actions artifact endpoint")
    origin = artifact.get("workflow_run", {})
    require(origin.get("id") == run_id and origin.get("head_sha") == sha and
            origin.get("head_branch") == "main" and origin.get("repository_id") == repository_id and
            origin.get("head_repository_id") == repository_id, "Actions artifact origin mismatch")
    require(start <= timestamp(artifact.get("created_at")) <= end,
            "Actions artifact was not created during original job")
    checksum = artifact.get("digest", "")
    require(isinstance(checksum, str) and checksum.startswith("sha256:") and digest(checksum[7:]),
            "Actions artifact digest missing or invalid")
    require(type(artifact.get("expired")) is bool, "unknown Actions artifact retention state")
    if artifact["expired"]:
        return False
    archive = api(f"{API}/actions/artifacts/{identifier}/zip", binary="actions")
    require(positive(artifact.get("size_in_bytes")) and len(archive) == artifact["size_in_bytes"] and
            hashlib.sha256(archive).hexdigest() == checksum[7:], "Actions archive size/digest mismatch")
    with zipfile.ZipFile(io.BytesIO(archive)) as zipped:
        entries = zipped.infolist()
        require(len(entries) == 1, "Actions archive must contain exactly one proof entry")
        entry = entries[0]
        mode = stat.S_IFMT(entry.external_attr >> 16)
        # Some ZIP producers omit Unix type bits. DOS directory and every
        # explicit non-regular Unix type are refused; no extraction takes place.
        require(entry.filename == entry.orig_filename == "fresh-acceptance-arm64.json" and not entry.is_dir() and
                not entry.external_attr & 0x10 and mode in (0, stat.S_IFREG) and
                not entry.flag_bits & 1 and entry.file_size <= PROOF_LIMIT,
                "unsafe or foreign Actions archive entry")
        require(zipped.read(entry) == proof, "Actions artifact differs from exact native proof")
    return True


def retained_proof(run_id, attempt, end, proof, body_sha):
    # Read candidate custody directly: no discovery, cache, CLI writes or trust
    # in a local similarly named receipt. The same sealed identity is read twice.
    def release_identity():
        release = api(f"{API}/releases/425202908")
        require(release.get("id") == 425202908 and release.get("draft") is True and
                release.get("tag_name") == "apple-candidate-v0.0.1" and
                release.get("target_commitish") == SOURCE and isinstance(release.get("body"), str) and
                hashlib.sha256(release["body"].encode()).hexdigest() == body_sha,
                "retained candidate identity/seal mismatch")

    release_identity()
    # Release inventories are arrays, unlike Actions inventories.
    assets, ids, names = [], set(), set()
    for page in range(1, 21):
        items = api(f"{API}/releases/425202908/assets?per_page=100&page={page}")
        require(isinstance(items, list) and len(items) <= 100, "invalid retained asset inventory")
        for item in items:
            identifier, name = item.get("id"), item.get("name")
            require(positive(identifier) and identifier not in ids and isinstance(name, str) and name not in names,
                    "duplicate retained asset ID/name")
            ids.add(identifier)
            names.add(name)
        assets.extend(items)
        if len(items) < 100:
            break
    else:
        raise Hold("retained asset pagination limit exceeded")
    name = f"fresh-acceptance-arm64-{run_id}-{attempt}.json"
    matches = [item for item in assets if item.get("name") == name]
    require(len(matches) == 1, "exact original retained fresh proof missing or ambiguous")
    asset = matches[0]
    require(asset.get("state") == "uploaded" and asset.get("size") == len(proof) and
            asset.get("url") == f"{URL}/releases/assets/{asset['id']}" and
            timestamp(asset.get("created_at")) >= end, "retained fresh asset binding mismatch")
    if asset.get("digest") is not None:
        require(asset["digest"] == "sha256:" + hashlib.sha256(proof).hexdigest(),
                "retained fresh asset digest mismatch")
    release_identity()
    remote = api(f"{API}/releases/assets/{asset['id']}", binary="release")
    require(remote == proof, "retained candidate fresh bytes differ from exact native proof")


def main():
    expected_env = {"GH_REPO": REPO, "NOTARY_RELEASE_ID": "425202908",
                    "NOTARY_RELEASE_SOURCE_SHA": SOURCE, "NOTARY_STATE_TAG": "apple-candidate-v0.0.1"}
    require(all(os.environ.get(key) == value for key, value in expected_env.items()),
            "exact candidate/repository environment pins required")
    body_sha = os.environ.get("NOTARY_RELEASE_BODY_SHA256")
    require(digest(body_sha), "candidate body SHA256 pin required")
    directory = Path(sys.argv[1])
    require(directory.is_dir() and not directory.is_symlink(), "unsafe evidence directory")
    path = directory / "acceptance.json"
    require(not path.is_symlink() and path.is_file() and path.stat().st_size <= PROOF_LIMIT,
            "native acceptance proof missing or unsafe")
    proof = path.read_bytes()
    receipt = parse(proof)
    expected = {"source_commit": SOURCE, "version": "0.0.1", "architecture": "arm64",
                "ci_run_id": "40470756538", "ci_run_attempt": "1",
                "acceptance_exception_parent_sha256": PARENT}
    require(all(receipt.get(key) == value for key, value in expected.items()) and
            type(receipt.get("acceptance_attempts")) is int and receipt["acceptance_attempts"] == 4 and
            receipt.get("credential_free_runner") is True and
            digest(receipt.get("zip_sha256")) and digest(receipt.get("dmg_sha256")),
            "not the exact exceptional ARM native proof")
    for key in ("verification_run_id", "verification_run_attempt"):
        require(isinstance(receipt.get(key), str) and re.fullmatch("[1-9][0-9]{0,19}", receipt[key]),
                "missing or invalid original verification identity")
    sha = receipt.get("verification_workflow_sha")
    require(digest(sha, 40), "original verification_workflow_sha is required")
    run_id, attempt = int(receipt["verification_run_id"]), int(receipt["verification_run_attempt"])
    require(run_id != 40470756538, "verification identity cannot be the original credential-free build")
    repository_id, start, end = verify_job(run_id, attempt, sha)
    name = f"macos-fresh-acceptance-arm64-{attempt}"
    artifacts = listing(f"{API}/actions/runs/{run_id}/artifacts?name={name}", "artifacts")
    require(all(item.get("name") == name for item in artifacts) and len(artifacts) <= 1,
            "fresh Actions artifact missing exact name or ambiguous")
    if artifacts and verify_archive(artifacts[0], run_id, sha, repository_id, start, end, proof):
        source = "original Actions artifact"
    else:
        retained_proof(run_id, attempt, end, proof, body_sha)
        source = "sealed candidate fresh receipt"
    # No success survives concurrent local proof replacement during reads.
    require(not path.is_symlink() and path.read_bytes() == proof, "native proof changed during completion audit")
    print(f"PASS: original fresh ARM job completed/success; exact proof via {source}; no qualification performed")


try:
    main()
except Hold as error:
    print("apple-fresh-completion: HOLD: " + str(error), file=sys.stderr)
    sys.exit(1)
except (OSError, ValueError, TypeError, KeyError, AttributeError, OverflowError,
        zipfile.BadZipFile, zlib.error, EOFError, RuntimeError, NotImplementedError, subprocess.TimeoutExpired):
    # No raw provider output, URLs, proof bytes, or parser exception text.
    print("apple-fresh-completion: HOLD: malformed, unavailable or timed-out completion evidence", file=sys.stderr)
    sys.exit(1)
PY
