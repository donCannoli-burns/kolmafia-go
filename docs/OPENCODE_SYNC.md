# OpenCode / ChatGPT GitHub synchronization

This document defines a **source synchronization workflow**, not permission to deploy or control a live KoLmafia runtime.

## Canonical shared GitHub state

- Repository: https://github.com/donCannoli-burns/kolmafia-go
- Stable branch: `main`
- Read-only workbench source: `dev/kol-workbench-v0.1/`
- Unresolved release gate: https://github.com/donCannoli-burns/kolmafia-go/issues/1
- Offline baseline evidence: [ISOLATED_BASELINE.md](ISOLATED_BASELINE.md)

`main` includes the source-archive-exact module README and executable mode `100755` for `dev/kol-workbench-v0.1/build.sh`. The original source/README bytes were not edited to change permissions.

## Native clone boundary

The original development and verifier root is:

`/home/sticky-ricky/repos/donCannoli-burns/don_kolmafia/Mafia-Builds/kolmafia-go`

**Do not** create a Git clone inside that directory. Do not reset or overwrite it. The inherited home-level Git repository is not a substitute for a dedicated project clone.

Proposed isolated **sibling** clone destination (requires agent-native approval and an absent target path):

`/home/sticky-ricky/repos/donCannoli-burns/kolmafia-go-github`

The agent should request/obtain its installed native filesystem/network approval for `git clone`. A chat instruction alone does not bypass installed tool permissions. Verify the clone destination is absent and outside the original project before cloning. If unavailable, stop and report the permission blocker; the GitHub connector can still review and edit repository files independently.

## Safe branch/PR lifecycle

1. Verify the GitHub `main` HEAD and compare source hashes with the frozen original.
2. Create a feature branch for a concrete, bounded repair; do not reuse or rewrite `main`.
3. Run tests and record exact commands, outputs and provenance for the proposed changes.
4. Push the feature branch only through an authorized Git operation.
5. Open a **draft PR** targeting `main` and link Issue #1. Do not merge until the missing checks and independent review are complete.
6. Leave the original project root, sibling loopback gateway, Docker containers, KoLmafia data, and game state untouched.

A GitHub issue tracks unresolved work; a draft PR represents an actual diff. Neither means deployment approval.

## Current scope limitation

The published Go source is an **offline-baselined, read-only development checkpoint**. No approval has been granted for runtime gateway takeover, Docker lifecycle commands, ASH/gCLI actions, mounts, deployments, or mutating API endpoints.
