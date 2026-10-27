#!/usr/bin/env bash
# Fixtures only: every provider command is intercepted and checked as exact GET.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
exec python3 - "$ROOT" "$@" <<'PY'
import copy
import hashlib
import io
import json
import os
from pathlib import Path
import shutil
import stat
import subprocess
import sys
import tempfile
import unittest
import zipfile

ROOT = Path(sys.argv[1])
HELPER = ROOT / "scripts/verify-apple-fresh-completion.sh"
EMITTING = len(sys.argv) in (5, 6) and sys.argv[2] == "--emit-fixture"
EXTERNAL = json.loads(Path(sys.argv[3]).read_bytes()) if EMITTING else {}
REPO = sys.argv[5] if EMITTING and len(sys.argv) == 6 else "Dunky13/takedown"
API = "repos/" + REPO
URL = "https://api.github.com/" + API
SOURCE = EXTERNAL.get("source_commit", "ef82381d509aac8b5b02da0f784a8176b1c5b9c0")
SHA = EXTERNAL.get("verification_workflow_sha", "a" * 40)
PARENT = "27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe"
RUN = int(EXTERNAL.get("verification_run_id", "777"))
ATTEMPT = int(EXTERNAL.get("verification_run_attempt", "2"))
RUN_API = f"{API}/actions/runs/{RUN}"
RUN_GET = RUN_API + f"/attempts/{ATTEMPT}"
JOBS = RUN_GET + "/jobs?per_page=100&page=1"
ARTIFACTS = RUN_API + f"/artifacts?name=macos-fresh-acceptance-arm64-{ATTEMPT}&per_page=100&page=1"
ZIP_GET = API + "/actions/artifacts/91/zip"
RELEASE_GET = API + "/releases/425202908"
ASSETS = RELEASE_GET + "/assets?per_page=100&page=1"
ASSET_GET = API + "/releases/assets/92"
CRITICAL = [
    "Qualify pinned automation before fresh acceptance",
    "Set up owned Apple ticket transport preference",
    "Install and assess exact final bytes on fresh native runner",
    "Verify owned cleanup before fresh receipt promotion",
    "Promote fresh receipt after verified transport cleanup",
]
MOCK = r'''#!/usr/bin/env python3
import json
import os
from pathlib import Path
import sys
base = Path(os.environ["FIXTURE"])
args = sys.argv[1:]
with (base / "calls.jsonl").open("a") as stream:
    stream.write(json.dumps(args) + "\n")
prefix = ["api", "--hostname", "github.com", "--method", "GET",
          "-H", "X-GitHub-Api-Version: 2022-11-28"]
assert args[:len(prefix)] == prefix, "provider mutation or foreign host"
endpoint = args[-1]
suffix = args[len(prefix):-1]
if endpoint.endswith("/zip"):
    assert suffix == ["--allow-escape-sequences", "-H", "Accept: application/vnd.github+json"]
elif "/releases/assets/" in endpoint:
    assert suffix == ["--allow-escape-sequences", "-H", "Accept: application/octet-stream"]
else:
    assert suffix == ["-H", "Accept: application/vnd.github+json"]
assert endpoint.startswith("repos/Dunky13/takedown/"), "unbound endpoint"
assert "GH_DEBUG" not in os.environ and "GH_FORCE_TTY" not in os.environ
manifest = json.loads((base / "responses.json").read_text())
if endpoint not in manifest or manifest[endpoint] == "FAIL":
    print("private SAS https://private.invalid/?sig=DO_NOT_LOG", file=sys.stderr)
    sys.exit(1)
sys.stdout.buffer.write((base / manifest[endpoint]).read_bytes())
'''


def raw_json(value):
    return (json.dumps(value, indent=2) + "\n").encode()


class CompletionFixtures(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="takedown-completion-test-")
        self.addCleanup(self.temporary.cleanup)
        self.base = Path(self.temporary.name)
        self.evidence = self.base / "evidence"
        self.evidence.mkdir()
        self.bin = self.base / "bin"
        self.bin.mkdir()
        (self.bin / "gh").write_text(MOCK.replace("repos/Dunky13/takedown/", "repos/" + REPO + "/"))
        (self.bin / "gh").chmod(0o700)
        # A mistaken native/signing/Apple/PF call is observable and always fails.
        for command in ("xcrun", "codesign", "spctl", "hdiutil", "pfctl", "security", "stapler", "notarytool"):
            (self.bin / command).write_text('#!/bin/sh\necho forbidden >> "$FIXTURE/forbidden"\nexit 99\n')
            (self.bin / command).chmod(0o700)
        self.receipt = {
            "source_commit": SOURCE, "version": "0.0.1", "architecture": "arm64",
            "ci_run_id": "40470756538", "ci_run_attempt": "1",
            "acceptance_attempts": 4, "acceptance_exception_parent_sha256": PARENT,
            "credential_free_runner": True, "verification_run_id": str(RUN),
            "verification_run_attempt": str(ATTEMPT), "verification_workflow_sha": SHA,
            "zip_sha256": "b" * 64, "dmg_sha256": "c" * 64,
        }
        self.proof = raw_json(self.receipt)
        (self.evidence / "acceptance.json").write_bytes(self.proof)
        (self.evidence / "acceptance-state.json").write_bytes(raw_json({
            "state": "verified", "attempts": 4, "acceptance_exception_parent_sha256": PARENT}))
        self.environment = os.environ.copy()
        self.environment.update({
            "PATH": str(self.bin) + os.pathsep + os.environ["PATH"], "FIXTURE": str(self.base),
            "GH_REPO": REPO, "NOTARY_RELEASE_ID": "425202908", "NOTARY_RELEASE_SOURCE_SHA": SOURCE,
            "NOTARY_STATE_TAG": "apple-candidate-v0.0.1",
            "NOTARY_RELEASE_BODY_SHA256": hashlib.sha256(b"sealed candidate body").hexdigest(),
            "GH_DEBUG": "api", "GH_FORCE_TTY": "1", "GH_TOKEN": "fixture-only",
        })
        self.responses = {}
        self.run = {
            "id": RUN, "run_attempt": ATTEMPT, "head_sha": SHA, "head_branch": "main",
            "event": "workflow_dispatch", "pull_requests": [], "workflow_id": 12,
            "status": "in_progress", "conclusion": None, "path": ".github/workflows/release.yml",
            "url": f"{URL}/actions/runs/{RUN}",
            "repository": {"id": 10, "full_name": REPO, "fork": False},
            "head_repository": {"id": 10, "full_name": REPO, "fork": False},
            "referenced_workflows": [{"path": REPO + "/.github/workflows/release-macos.yml@" + SHA,
                                      "sha": SHA, "ref": "refs/heads/main"}],
        }
        self.job = {
            "id": 31, "name": "macos / acceptance (macos-15, arm64)", "run_id": RUN,
            "run_attempt": ATTEMPT, "head_sha": SHA, "head_branch": "main",
            "run_url": f"{URL}/actions/runs/{RUN}", "url": f"{URL}/actions/jobs/31",
            "labels": ["macos-15"], "status": "completed", "conclusion": "success",
            "started_at": "2026-10-27T10:00:00Z", "completed_at": "2026-10-27T10:10:00Z",
            "steps": [{"name": name, "number": number, "status": "completed", "conclusion": "success"}
                      for number, name in enumerate(CRITICAL, 1)],
        }
        self.set_json(RUN_GET, self.run)
        self.set_json(API + "/actions/workflows/12", {"id": 12, "path": ".github/workflows/release.yml"})
        self.set_json(API + "/branches/main", {"name": "main", "protected": True, "commit": {"sha": SHA}})
        self.set_json(API + f"/compare/{SHA}...{SHA}?per_page=1", {
            "status": "identical", "base_commit": {"sha": SHA}, "merge_base_commit": {"sha": SHA}})
        self.jobs([self.job])
        self.artifact = {
            "id": 91, "name": f"macos-fresh-acceptance-arm64-{ATTEMPT}", "expired": False,
            "created_at": "2026-10-27T10:09:00Z", "url": URL + "/actions/artifacts/91",
            "archive_download_url": URL + "/actions/artifacts/91/zip",
            "workflow_run": {"id": RUN, "head_sha": SHA, "head_branch": "main",
                             "repository_id": 10, "head_repository_id": 10},
        }
        self.archive()
        self.release = {"id": 425202908, "draft": True, "tag_name": "apple-candidate-v0.0.1",
                        "target_commitish": SOURCE, "body": "sealed candidate body"}
        self.set_json(RELEASE_GET, self.release)
        self.asset = {"id": 92, "name": f"fresh-acceptance-arm64-{RUN}-{ATTEMPT}.json",
                      "state": "uploaded", "size": len(self.proof), "created_at": "2026-10-27T10:11:00Z",
                      "url": URL + "/releases/assets/92",
                      "digest": "sha256:" + hashlib.sha256(self.proof).hexdigest()}
        self.set_json(ASSETS, [self.asset])
        self.set_raw(ASSET_GET, self.proof)

    def set_raw(self, endpoint, raw):
        filename = hashlib.sha256(endpoint.encode()).hexdigest()
        (self.base / filename).write_bytes(raw)
        self.responses[endpoint] = filename

    def set_json(self, endpoint, value):
        self.set_raw(endpoint, raw_json(value))

    def jobs(self, jobs):
        self.set_json(JOBS, {"total_count": len(jobs), "jobs": jobs})

    def artifacts(self, artifacts):
        self.set_json(ARTIFACTS, {"total_count": len(artifacts), "artifacts": artifacts})

    def archive(self, entries=None, raw=None):
        if raw is None:
            buffer = io.BytesIO()
            with zipfile.ZipFile(buffer, "w") as archive:
                for entry, content in entries or [("fresh-acceptance-arm64.json", self.proof)]:
                    archive.writestr(entry, content)
            raw = buffer.getvalue()
        self.artifact.update({"size_in_bytes": len(raw), "digest": "sha256:" + hashlib.sha256(raw).hexdigest()})
        self.set_raw(ZIP_GET, raw)
        self.artifacts([self.artifact])

    def check(self, success, contains=None):
        (self.base / "responses.json").write_bytes(raw_json(self.responses))
        def snapshot():
            return {str(path.relative_to(self.evidence)): (path.is_symlink(),
                    os.readlink(path) if path.is_symlink() else path.read_bytes())
                    for path in self.evidence.rglob("*") if path.is_file() or path.is_symlink()}
        before = snapshot()
        result = subprocess.run(["/bin/bash", str(HELPER), str(self.evidence)], env=self.environment,
                                capture_output=True, text=True, timeout=15, check=False)
        self.assertEqual(result.returncode == 0, success, result.stdout + result.stderr)
        self.assertEqual(before, snapshot(), "helper mutated retained evidence/counter")
        self.assertFalse((self.base / "forbidden").exists(), "native/Apple/PF/signing call")
        self.assertNotIn("DO_NOT_LOG", result.stdout + result.stderr)
        self.assertNotIn("private.invalid", result.stdout + result.stderr)
        if contains:
            self.assertIn(contains, result.stdout + result.stderr)
        calls_path = self.base / "calls.jsonl"
        self.calls = [json.loads(line) for line in calls_path.read_text().splitlines()] if calls_path.exists() else []
        for call in self.calls:
            self.assertEqual(call[:5], ["api", "--hostname", "github.com", "--method", "GET"])
            self.assertNotIn("--input", call)
            self.assertNotIn("--paginate", call)
            self.assertNotIn("--field", call)
            self.assertNotIn("--raw-field", call)
            self.assertNotIn(f"{RUN_API}/jobs", call[-1])
        return result

    def test_exact_success_while_workflow_still_running(self):
        self.check(True, "original Actions artifact")
        self.assertIn(RUN_GET, [call[-1] for call in self.calls])
        self.assertIn(JOBS, [call[-1] for call in self.calls])
        self.assertNotIn(RELEASE_GET, [call[-1] for call in self.calls])

    def test_publication_retry_uses_original_job_no_new_attempt(self):
        self.run.update(status="completed", conclusion="failure", run_attempt=ATTEMPT)
        self.set_json(RUN_GET, self.run)
        self.environment.update(GITHUB_RUN_ID="9999", GITHUB_RUN_ATTEMPT="7", RELEASE_AUTOMATION_SHA="d" * 40)
        self.check(True)
        self.assertFalse(any("9999" in call[-1] for call in self.calls))

    def test_checkpoint_late_commit_cancellation_restored_proof_refused(self):
        # Immutable verified4 checkpoint commits remotely while cancellation
        # removes local success. Recreate a later wrapper/readiness restore.
        durable = self.base / "durable-verified4"
        shutil.copytree(self.evidence, durable)
        (self.evidence / "acceptance.json").unlink()
        shutil.rmtree(self.evidence)
        shutil.copytree(durable, self.evidence)
        # This is the old checkpoint-only success predicate, which remains
        # true after the immutable late commit and would permit replay.
        self.assertEqual(json.loads((self.evidence / "acceptance-state.json").read_text())["state"], "verified")
        self.assertEqual((self.evidence / "acceptance.json").read_bytes(), self.proof)
        self.job["conclusion"] = "cancelled"
        self.jobs([self.job])
        self.check(False, "not completed/success")
        self.assertFalse(any(call[-1] in (ARTIFACTS, ASSETS, ZIP_GET, ASSET_GET) for call in self.calls))
        self.assertEqual(json.loads((self.evidence / "acceptance-state.json").read_text())["attempts"], 4)

    def test_failed_pending_skipped_and_missing_original_job(self):
        for status, conclusion in [("completed", "failure"), ("in_progress", None),
                                   ("queued", None), ("completed", "skipped"), ("completed", "timed_out")]:
            with self.subTest(status=status, conclusion=conclusion):
                job = dict(self.job, status=status, conclusion=conclusion)
                self.jobs([job])
                self.check(False)
        self.jobs([])
        self.check(False)

    def test_green_intel_or_other_job_cannot_replace_arm(self):
        for name in ("macos / acceptance (macos-15-intel, amd64)", "macos / sign (macos-15, arm64)"):
            with self.subTest(name=name):
                self.jobs([dict(self.job, name=name)])
                self.check(False)

    def test_duplicate_jobs_and_ids_refused(self):
        for jobs in ([self.job, self.job], [self.job, dict(self.job, id=32)]):
            self.jobs(jobs)
            self.check(False)

    def test_wrong_job_identity_and_architecture(self):
        for key, value in {"run_id": 778, "run_attempt": 3, "head_sha": "d" * 40,
                           "head_branch": "other", "labels": ["macos-15-intel"],
                           "run_url": "https://foreign.invalid/run", "url": URL + "/actions/jobs/32"}.items():
            with self.subTest(key=key):
                self.jobs([dict(self.job, **{key: value})])
                self.check(False)

    def test_original_attempt_no_latest_fallback(self):
        self.responses[RUN_GET] = "FAIL"
        self.set_json(RUN_API, dict(self.run, run_attempt=3, conclusion="success"))
        self.check(False)
        self.assertNotIn(RUN_API, [call[-1] for call in self.calls])

    def test_wrong_run_attempt_sha_workflow_repo_event(self):
        changes = [{"id": 778}, {"run_attempt": 3}, {"head_sha": "d" * 40}, {"head_branch": "other"},
                   {"status": "unknown"}, {"status": "completed", "conclusion": None},
                   {"path": ".github/workflows/ci.yml"}, {"path": ".github/workflows/release.yml@other"},
                   {"event": "pull_request"}, {"pull_requests": [{"id": 1}]},
                   {"repository": {"full_name": "foreign/repo", "id": 10, "fork": False}},
                   {"head_repository": {"full_name": REPO, "id": 11, "fork": False}},
                   {"head_repository": {"full_name": REPO, "id": 10, "fork": True}}]
        for change in changes:
            with self.subTest(change=change):
                self.set_json(RUN_GET, dict(self.run, **change))
                self.check(False)

    def test_wrong_workflow_id_or_reusable_sha(self):
        self.set_json(API + "/actions/workflows/12", {"id": 13, "path": ".github/workflows/release.yml"})
        self.check(False)
        self.set_json(API + "/actions/workflows/12", {"id": 12, "path": ".github/workflows/ci.yml"})
        self.check(False)
        self.set_json(API + "/actions/workflows/12", {"id": 12, "path": ".github/workflows/release.yml"})
        for references in ([], self.run["referenced_workflows"] * 2,
                           [dict(self.run["referenced_workflows"][0], sha="e" * 40)],
                           [dict(self.run["referenced_workflows"][0], ref="refs/heads/other")],
                           [dict(self.run["referenced_workflows"][0], ref=None)]):
            self.set_json(RUN_GET, dict(self.run, referenced_workflows=references))
            self.check(False)

    def test_protected_main_and_ancestry_required(self):
        self.set_json(API + "/branches/main", {"name": "main", "protected": False, "commit": {"sha": SHA}})
        self.check(False)
        self.set_json(API + "/branches/main", {"name": "main", "protected": True, "commit": {"sha": SHA}})
        for change in ({"status": "diverged"}, {"merge_base_commit": {"sha": "d" * 40}}):
            response = {"status": "identical", "base_commit": {"sha": SHA}, "merge_base_commit": {"sha": SHA}}
            response.update(change)
            self.set_json(API + f"/compare/{SHA}...{SHA}?per_page=1", response)
            self.check(False)

    def test_each_critical_step_must_succeed_once_in_order(self):
        for index in range(len(CRITICAL)):
            for conclusion in ("cancelled", "failure", "skipped", None):
                with self.subTest(index=index, conclusion=conclusion):
                    job = copy.deepcopy(self.job)
                    job["steps"][index]["conclusion"] = conclusion
                    self.jobs([job])
                    self.check(False)
        for steps in (self.job["steps"][:-1], self.job["steps"] + [dict(self.job["steps"][0], number=6)]):
            self.jobs([dict(self.job, steps=steps)])
            self.check(False)
        steps = copy.deepcopy(self.job["steps"])
        steps[0]["number"], steps[1]["number"] = 2, 1
        self.jobs([dict(self.job, steps=steps)])
        self.check(False)

    def test_job_pagination_includes_exact_match_on_second_page(self):
        others = [{"id": 1000 + index, "name": "other"} for index in range(100)]
        self.set_json(JOBS, {"total_count": 101, "jobs": others})
        self.set_json(JOBS.replace("&page=1", "&page=2"), {"total_count": 101, "jobs": [self.job]})
        self.check(True)

    def test_pagination_duplicates_and_changed_counts_hold(self):
        others = [{"id": 1000 + index, "name": "other"} for index in range(99)] + [self.job]
        self.set_json(JOBS, {"total_count": 101, "jobs": others})
        for page in ({"total_count": 101, "jobs": [self.job]}, {"total_count": 102, "jobs": [dict(self.job, id=32)]}):
            self.set_json(JOBS.replace("&page=1", "&page=2"), page)
            self.check(False)
        self.set_json(JOBS, {"total_count": 2, "jobs": [self.job]})
        self.check(False)

    def test_missing_native_provenance_and_foreign_scope(self):
        changes = {"verification_workflow_sha": None, "verification_run_id": "40470756538",
                   "verification_run_attempt": "3", "architecture": "amd64", "source_commit": "f" * 40,
                   "ci_run_id": "40470756539", "ci_run_attempt": "2", "version": "0.0.2",
                   "acceptance_attempts": 3, "credential_free_runner": False,
                   "acceptance_exception_parent_sha256": "d" * 64, "zip_sha256": "bad"}
        for key, value in changes.items():
            with self.subTest(key=key):
                (self.evidence / "acceptance.json").write_bytes(raw_json(dict(self.receipt, **{key: value})))
                self.check(False)

    def test_wrong_environment_pins_hold_without_provider_calls(self):
        for key in ("GH_REPO", "NOTARY_RELEASE_ID", "NOTARY_RELEASE_SOURCE_SHA", "NOTARY_STATE_TAG", "NOTARY_RELEASE_BODY_SHA256"):
            with self.subTest(key=key):
                previous = self.environment[key]
                self.environment[key] = "wrong"
                self.check(False)
                self.environment[key] = previous
        self.assertEqual(self.calls, [])

    def test_missing_symlink_and_malformed_native_proof(self):
        path = self.evidence / "acceptance.json"
        path.unlink()
        self.check(False)
        target = self.base / "native.json"
        target.write_bytes(self.proof)
        path.symlink_to(target)
        self.check(False)
        path.unlink()
        for raw in (b"{", b"[]", b'{"a":1,"a":2}'):
            path.write_bytes(raw)
            self.check(False)

    def test_duplicate_artifacts_and_wrong_names_hold(self):
        for artifacts in ([self.artifact, self.artifact], [self.artifact, dict(self.artifact, id=93)],
                          [dict(self.artifact, name="macos-fresh-acceptance-amd64-2")],
                          [dict(self.artifact, name="macos-fresh-acceptance-arm64-3")]):
            self.artifacts(artifacts)
            self.check(False)

    def test_artifact_pagination_ambiguity_cannot_hide_on_second_page(self):
        artifacts = [dict(self.artifact, id=1000 + index) for index in range(100)]
        self.set_json(ARTIFACTS, {"total_count": 101, "artifacts": artifacts})
        self.set_json(ARTIFACTS.replace("&page=1", "&page=2"), {"total_count": 101, "artifacts": [self.artifact]})
        self.check(False)

    def test_artifact_source_run_sha_repo_branch_and_time_required(self):
        for key, value in {"id": 778, "head_sha": "d" * 40, "repository_id": 11,
                           "head_repository_id": 11, "head_branch": "other"}.items():
            with self.subTest(key=key):
                origin = dict(self.artifact["workflow_run"], **{key: value})
                self.artifacts([dict(self.artifact, workflow_run=origin)])
                self.check(False)
        for change in ({"created_at": "2026-10-27T10:11:00Z"}, {"url": "https://foreign.invalid"},
                       {"archive_download_url": "https://foreign.invalid"}, {"expired": None},
                       {"digest": None}, {"digest": "sha256:" + "d" * 64}, {"size_in_bytes": 1}):
            self.artifacts([dict(self.artifact, **change)])
            self.check(False)

    def test_wrong_outer_archive_bytes_and_digest_never_fall_back(self):
        self.set_raw(ZIP_GET, b"corrupt")
        self.check(False, "digest mismatch")
        self.assertNotIn(RELEASE_GET, [call[-1] for call in self.calls])

    def test_archive_malformed_foreign_extra_duplicate_symlink_and_directory(self):
        link = zipfile.ZipInfo("fresh-acceptance-arm64.json")
        link.create_system = 3
        link.external_attr = (stat.S_IFLNK | 0o777) << 16
        folder = zipfile.ZipInfo("fresh-acceptance-arm64.json")
        folder.create_system = 3
        folder.external_attr = ((stat.S_IFDIR | 0o755) << 16) | 0x10
        cases = [[], [("other.json", self.proof)], [("../fresh-acceptance-arm64.json", self.proof)],
                 [("/fresh-acceptance-arm64.json", self.proof)], [("fresh-acceptance-arm64.json/", self.proof)],
                 [("fresh-acceptance-arm64.json", self.proof), ("extra", b"extra")],
                 [(link, b"target")], [(folder, b"")], [("fresh-acceptance-arm64.json", b"foreign proof")],
                 [("fresh-acceptance-arm64.json", b" " * (1024 * 1024 + 1))]]
        for entries in cases:
            with self.subTest(entries=[str(item[0]) for item in entries]):
                if not entries:
                    self.archive(raw=b"not a ZIP")
                else:
                    self.archive(entries)
                self.check(False)
        # Duplicate names are a separate central-directory ambiguity.
        import warnings
        with warnings.catch_warnings(record=True):
            self.archive([("fresh-acceptance-arm64.json", self.proof)] * 2)
        self.check(False)

    def test_empty_archive_encryption_and_crc_corruption_hold(self):
        buffer = io.BytesIO()
        with zipfile.ZipFile(buffer, "w"):
            pass
        self.archive(raw=buffer.getvalue())
        self.check(False)
        self.archive()
        valid = (self.base / self.responses[ZIP_GET]).read_bytes()
        encrypted = bytearray(valid)
        # Mark the sole entry encrypted in both local and central headers.
        encrypted[6] |= 1
        central = encrypted.index(b"PK\x01\x02")
        encrypted[central + 8] |= 1
        self.archive(raw=bytes(encrypted))
        self.check(False)
        corrupt = bytearray(valid)
        # Outer digest is recomputed, leaving a failing inner CRC check.
        corrupt[valid.index(self.proof)] ^= 1
        self.archive(raw=bytes(corrupt))
        self.check(False)

    def test_missing_artifact_and_missing_retained_proof_hold(self):
        self.artifacts([])
        self.set_json(ASSETS, [])
        self.check(False, "retained fresh proof missing")

    def test_retained_asset_pagination_and_old_automation_on_current_main(self):
        self.artifacts([])
        others = [dict(self.asset, id=1000 + index, name=f"unrelated-{index}") for index in range(100)]
        self.set_json(ASSETS, others)
        self.set_json(ASSETS.replace("&page=1", "&page=2"), [self.asset])
        new_main = "d" * 40
        self.set_json(API + "/branches/main", {"name": "main", "protected": True, "commit": {"sha": new_main}})
        self.set_json(API + f"/compare/{SHA}...{new_main}?per_page=1", {
            "status": "ahead", "base_commit": {"sha": SHA}, "merge_base_commit": {"sha": SHA}})
        self.check(True, "sealed candidate fresh receipt")

    def test_same_json_different_bytes_refused(self):
        self.archive([("fresh-acceptance-arm64.json", json.dumps(self.receipt).encode())])
        self.check(False, "exact native proof")

    def test_missing_or_expired_actions_artifact_uses_bound_durable_receipt(self):
        self.artifacts([])
        self.check(True, "sealed candidate fresh receipt")
        self.assertIn(ASSET_GET, [call[-1] for call in self.calls])
        self.artifacts([dict(self.artifact, expired=True)])
        self.check(True, "sealed candidate fresh receipt")

    def test_durable_fallback_keeps_original_success_requirement(self):
        self.artifacts([])
        self.job["conclusion"] = "cancelled"
        self.jobs([self.job])
        self.check(False)
        self.assertNotIn(RELEASE_GET, [call[-1] for call in self.calls])

    def test_unbound_missing_duplicate_foreign_retained_assets_hold(self):
        self.artifacts([])
        changes = [{"name": f"fresh-acceptance-arm64-999-{ATTEMPT}.json"},
                   {"name": f"fresh-acceptance-arm64-{RUN}-3.json"},
                   {"state": "starter"}, {"size": 1}, {"url": "https://foreign.invalid"},
                   {"created_at": "2026-10-27T10:09:00Z"}, {"digest": "sha256:" + "f" * 64}]
        for assets in ([], [self.asset, self.asset], [self.asset, dict(self.asset, id=93)]):
            self.set_json(ASSETS, assets)
            self.check(False)
        for change in changes:
            self.set_json(ASSETS, [dict(self.asset, **change)])
            self.check(False)
        self.set_json(ASSETS, [self.asset])
        self.set_raw(ASSET_GET, b"foreign proof")
        self.check(False)

    def test_candidate_exact_id_body_source_and_draft_required(self):
        self.artifacts([])
        for change in ({"id": 424598089}, {"body": "other"}, {"target_commitish": "d" * 40},
                       {"draft": False}, {"tag_name": "apple-candidate-v0.0.2"}):
            self.set_json(RELEASE_GET, dict(self.release, **change))
            self.check(False)

    def test_api_failures_never_fall_back_or_leak_sas_urls(self):
        baseline = dict(self.responses)
        for endpoint in (RUN_GET, API + "/actions/workflows/12", API + "/branches/main", JOBS, ARTIFACTS, ZIP_GET):
            with self.subTest(endpoint=endpoint):
                self.responses = dict(baseline)
                self.responses[endpoint] = "FAIL"
                self.check(False, "authoritative API read failed")
        self.responses = dict(baseline)
        self.artifacts([])
        baseline = dict(self.responses)
        for endpoint in (RELEASE_GET, ASSETS, ASSET_GET):
            self.responses = dict(baseline)
            self.responses[endpoint] = "FAIL"
            self.check(False)

    def test_malformed_api_and_duplicate_json_keys_hold(self):
        for raw in (b"{", b"[]", b'{"id":777,"id":778}', b"null", b'{"id":NaN}'):
            self.set_raw(RUN_GET, raw)
            self.check(False)


if EMITTING:
    # Shared private API fixture for real wrapper/publisher integration tests.
    # Production scripts have no fixture mode or completion-success override.
    fixture = CompletionFixtures()
    try:
        fixture.setUp()
        fixture.proof = Path(sys.argv[3]).read_bytes()
        fixture.archive()
        fixture.asset.update(size=len(fixture.proof), digest="sha256:" + hashlib.sha256(fixture.proof).hexdigest())
        fixture.set_json(ASSETS, [fixture.asset])
        fixture.set_raw(ASSET_GET, fixture.proof)
        (fixture.base / "responses.json").write_bytes(raw_json(fixture.responses))
        pins = {key: fixture.environment[key] for key in ("GH_REPO", "NOTARY_RELEASE_ID", "NOTARY_RELEASE_SOURCE_SHA", "NOTARY_STATE_TAG", "NOTARY_RELEASE_BODY_SHA256")}
        (fixture.base / "pins.json").write_bytes(raw_json(pins))
        shutil.copytree(fixture.base, Path(sys.argv[4]))
    finally:
        fixture.doCleanups()
else:
    unittest.main(argv=[sys.argv[0]], verbosity=2)
PY
