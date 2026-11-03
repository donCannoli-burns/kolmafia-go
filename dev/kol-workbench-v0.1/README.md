# KoL Workbench — v0.1 read-only starter

A local Go gateway and go-app v11 WebAssembly dashboard. This first milestone **only reads** the existence of five registered source/data directories and three tool binaries, then inspects the state of the existing `kolmafia-xpra-kolmafia-1` Docker container. It does not start or stop containers, run KoLmafia scripts, lint JSON, or write to user data.

## Known host facts (from 2026-11-01 terminal inventory)

- Go: 1.26.8 linux/amd64; Docker Engine: 29.7.2; Docker CLI: `/usr/bin/docker`.
- KoLmafia container: `kolmafia-xpra-kolmafia-1`, stopped (`Exited 143`), no published ports while stopped.
- Existing writable volume mapping: `~/repos/es-3581100/kolmafia-xpra/data` -> `/home/kolmafia/.kolmafia`.
- Existing Docker named volume: `kolmafia-xpra_jars` -> `/home/kolmafia/jars`.
- `jsonlint`: `~/.local/bin/jsonlint`; `fx`: `/usr/local/bin/fx`.
- Source/data paths are registered explicitly in `internal/host/status.go`, including both `data-of-loathing` paths and the ASH catalog index. Other proposed repositories are *not* yet verified.
- Docker did not report a rootless security option; do not assume rootless engine configuration.

## Build on your Linux machine

Prerequisites: Go 1.26+ and normal access to the Go module proxy to download `go-app/v11` once.

```sh
cd ~/repos/donCannoli-burns
# Extract this archive to a fresh kol-workbench-v0.1 directory here.
cd kol-workbench-v0.1
go mod tidy
go test ./internal/...
bash build.sh
./bin/kol-workbench
```

Open **http://127.0.0.1:18484/** and press `Refresh inventory`. Run `curl -fsS http://127.0.0.1:18484/api/v1/status` to inspect JSON. Because go-app resolves its web resources relative to the working directory, run the server from this project's root.

## Read-only API

- `GET /api/v1/health` -> `{ "status": "ok", "mode": "read-only", "schema_version": 1 }`
- `GET /api/v1/status` -> bounded inventory snapshot

The gateway rejects non-GET/HEAD methods and unexpected Host/Origin headers. It binds only to `127.0.0.1:18484` and sends no permissive CORS headers. **No authentication has been added**, so do not add any mutating actions before implementing user authentication, CSRF protection, and explicit operation authorization. Read-only mode is not a claim that host file paths are harmless to disclose to local users.

## Future milestones, not implemented

1. Authenticated local capability registry + bounded CLI validation jobs.
2. Build artifact staging and source-hash provenance.
3. Read-only game-data adapter for `data-of-loathing`.
4. SSE events with snapshot recovery; semantic SVG state model.
5. Explicit approvals for deployments and lifecycle actions using existing kol-playctl / Compose conventions.

## Safety contract

- This code calls **only** `docker inspect --format '{{json .State}}' kolmafia-xpra-kolmafia-1`.
- No arbitrary commands or paths are accepted from browser requests.
- Nothing is mounted or written into KoLmafia's data directory.
- A stopped KoLmafia container stays stopped.
- Only the five discovered repositories/indexes and three known tools are inspected.

## Verification note

The host API package has standard-library-only unit tests; go-app v11 needs Go 1.26+, while the packaging environment used to generate this archive has Go 1.23.2. Therefore the **complete go-app/WASM build remains to be validated on your Go 1.26.8 machine**. Treat this as a scaffold, not an already validated production executable.
