// Package host provides narrowly scoped, read-only discovery of approved host resources.
package host

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const ContainerName = "kolmafia-xpra-kolmafia-1"

type Resource struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Path    string `json:"path"`
	Present bool   `json:"present"`
	Kind    string `json:"kind"`
}

type Container struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Running   bool   `json:"running"`
	Status    string `json:"status"`
	ExitCode  int    `json:"exit_code"`
}

type Snapshot struct {
	SchemaVersion int        `json:"schema_version"`
	CollectedAt   time.Time  `json:"collected_at"`
	Mode          string     `json:"mode"`
	Repositories  []Resource `json:"repositories"`
	Tools         []Resource `json:"tools"`
	Container     Container  `json:"container"`
}

type Target struct{ ID, Label, Path, Kind string }

// These are known inventory candidates, not a broad permission to access arbitrary paths.
// Other repositories will be added only after their locations have been confirmed.
func DefaultRepositories(home string) []Target {
	root := filepath.Join(home, "repos")
	return []Target{
		{"kolmafia-xpra", "KoLmafia Xpra", filepath.Join(root, "es-3581100", "kolmafia-xpra"), "source"},
		{"kolmafia-ts", "KoLmafia TS declarations", filepath.Join(root, "donCannoli-burns", "don_kolmafia", "Mafia-Builds", "kolmafia-ts"), "source"},
		{"data-of-loathing-source", "Data of Loathing / source", filepath.Join(root, "kol-watersheep-source-stack", "sources", "data-of-loathing"), "source"},
		{"data-of-loathing-data", "Data of Loathing / data", filepath.Join(root, "kol-watersheep-source-stack", "data", "data-of-loathing"), "data"},
		{"master-ash-catalog", "Master ASH catalog index", filepath.Join(root, "kol-watersheep-source-stack", "indexes", "master-ash-catalog"), "index"},
	}
}

func DefaultTools(home string) []Target {
	return []Target{
		{"jsonlint", "jsonlint CLI", filepath.Join(home, ".local", "bin", "jsonlint"), "executable"},
		{"fx", "fx CLI", "/usr/local/bin/fx", "executable"},
		{"docker", "Docker CLI", "/usr/bin/docker", "executable"},
	}
}

func inspectResources(targets []Target) []Resource {
	out := make([]Resource, 0, len(targets))
	for _, t := range targets {
		present := false
		info, err := os.Stat(t.Path)
		if err == nil {
			if t.Kind == "executable" {
				present = !info.IsDir() && info.Mode()&0111 != 0
			} else {
				present = info.IsDir()
			}
		}
		out = append(out, Resource{ID: t.ID, Label: t.Label, Path: t.Path, Present: present, Kind: t.Kind})
	}
	return out
}

// InspectContainer issues exactly one Docker inspect request for the fixed, known container.
// It never starts, stops, creates, removes, or attaches to any container.
func InspectContainer(ctx context.Context, dockerPath string) Container {
	c := Container{Name: ContainerName, Status: "unavailable"}
	limited, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	command := exec.CommandContext(limited, dockerPath, "inspect", "--format", "{{json .State}}", ContainerName)
	output, err := command.Output()
	if err != nil {
		return c
	}
	var state struct {
		Status   string `json:"Status"`
		Running  bool   `json:"Running"`
		ExitCode int    `json:"ExitCode"`
	}
	if err := json.Unmarshal(output, &state); err != nil || strings.TrimSpace(state.Status) == "" {
		return c
	}
	c.Available = true
	c.Running = state.Running
	c.Status = state.Status
	c.ExitCode = state.ExitCode
	return c
}

func Collect(ctx context.Context, home string) (Snapshot, error) {
	return collect(ctx, home, inspectResources, time.Now, func(ctx context.Context) Container {
		return InspectContainer(ctx, "/usr/bin/docker")
	})
}

func collect(ctx context.Context, home string, inspect func([]Target) []Resource, now func() time.Time, container func(context.Context) Container) (Snapshot, error) {
	if home == "" || !filepath.IsAbs(home) {
		return Snapshot{}, errors.New("absolute home directory required")
	}
	repos := inspect(DefaultRepositories(home))
	tools := inspect(DefaultTools(home))
	return Snapshot{
		SchemaVersion: 1,
		CollectedAt:   now().UTC(),
		Mode:          "read-only",
		Repositories:  repos,
		Tools:         tools,
		Container:     container(ctx),
	}, nil
}
