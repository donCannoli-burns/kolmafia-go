# KoLmafia-Go / KoL Workbench

This repository publishes the **isolated, read-only Go + go-app/WASM baseline** from the KoLmafia-Go development workspace. It is a development checkpoint, **not** a release or live KoLmafia controller.

## Current implementation

Source is in [`dev/kol-workbench-v0.1/`](dev/kol-workbench-v0.1/). It provides a loopback Go HTTP gateway, a go-app v11 WebAssembly dashboard, fixed host resource inventory, and read-only inspection of a named Docker container. There are no browser-facing CLI execution, build, deployment, or container lifecycle routes. It has no authentication; do not expose it on a network or extend it with mutating routes without an authorization review.

## Verified scope

The project's November 2, 2026 isolated baseline records **four fake-backed host tests passing**, **native Go compilation**, **JS/WASM compilation**, **go vet**, and **format verification**. These were offline checks reported by the project verifier, not a claim of successful production or interactive runtime operation. The verified source was a frozen nine-file subset, with no listener or live game actions exercised.

## Open blockers

- Full application/runtime behavior, security, dependency and toolchain provenance remain unverified.
- No runtime, deployment, live KoLmafia mutation, or sibling gateway takeover has been authorized.
- Follow-up work needs bounded contracts and independent review before enabling new capabilities.

See the repository's open GitHub issues for the current work and release gates.

## Local development

On a host with Go 1.26+:

```sh
cd dev/kol-workbench-v0.1
go test ./internal/...
go vet ./...
GOOS=js GOARCH=wasm go build -o /tmp/workbench.wasm ./cmd/workbench
go build -o /tmp/workbench-native ./cmd/workbench
```

The commands above are for compilation/testing only. Running the resulting gateway or using Docker requires separate local safety review; the default server port can conflict with an existing sibling deployment.

## Source provenance

Published from the user-supplied `kolmafia-go.zip` snapshot (SHA-256 `b480d12b3b19d73fe7188f6c27980e09c07e50cce6049e73563bfeb4c0bc9b15`). The public tree intentionally omits caches, generated binaries, local audit/agent records, snapshots and private machine configuration.
