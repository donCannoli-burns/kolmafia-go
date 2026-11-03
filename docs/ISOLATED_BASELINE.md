# Isolated baseline — evidence summary (2026-11-02)

**Status:** `completed` within isolated baseline scope. **Overall runtime safety:** partial and unverified.

This public summary is derived from the project checkpoint and accompanying independent verification documentation in the uploaded November 2, 2026 snapshot. It is not an independent rerun of the Go commands by this publication process.

## Recorded checks

| Check | Reported outcome |
| --- | --- |
| Fake-backed host tests | PASS — 4 tests |
| Native compilation | PASS — output not executed |
| JS/WASM compilation | PASS — output not executed |
| `go vet` | PASS |
| Format verification | PASS |

Independent B2 review and B3 documentation were published and their hashes rechecked in the canonical return. The source was frozen during the bounded review. Full `go test ./...`, interactive browser runtime, listener behavior and deployment were not verified.

## Open medium-severity review findings

1. A running, unauthenticated loopback gateway can disclose local resource paths and trigger repeated fixed Docker inspection via the status route.
2. Test coverage omits status API, HEAD and allowed-origin cases, invalid home, nil callbacks, Docker failures/timeouts and third-party handler behavior.
3. Dependency/toolchain provenance, full application behavior and native/WASM interactive runtime remain unaudited.

## Safety boundary

No deployment, live KoLmafia action, container lifecycle operation, sibling gateway replacement or new privileged route has been approved. The dormant feature implementation plan is not released. Deployment should remain blocked until separately reviewed and authorized.

## Provenance

- Uploaded full workspace archive SHA-256: `b480d12b3b19d73fe7188f6c27980e09c07e50cce6049e73563bfeb4c0bc9b15`.
- Canonical B2/B3 publication checkpoint SHA-256: `ffa75c6189ef4175c7240464a541950ec45af0451fb3f31cca76d542596700c8`.
- The public source subtree contains eight exact-byte source files from `dev/kol-workbench-v0.1/`; checked against their original Git blob hashes upon publication.
- Raw local verifier records, build caches, agent runtime captures and host-specific archived reports were intentionally excluded from the public repository.

Treat hashes as provenance, not assurance of safety or correctness.
