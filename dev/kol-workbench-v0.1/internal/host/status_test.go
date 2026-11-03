package host

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCollectDoesNotRequireDocker(t *testing.T) {
	home := t.TempDir()
	wantTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	resourceCalls, containerCalls := 0, 0
	fakeResources := func(targets []Target) []Resource {
		resourceCalls++
		out := make([]Resource, 0, len(targets))
		for _, target := range targets {
			out = append(out, Resource{ID: target.ID, Label: target.Label, Path: target.Path, Kind: target.Kind})
		}
		return out
	}
	fakeContainer := func(context.Context) Container {
		containerCalls++
		return Container{Name: ContainerName, Status: "unavailable"}
	}
	got, err := collect(context.Background(), home, fakeResources, func() time.Time { return wantTime }, fakeContainer)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != "read-only" || got.SchemaVersion != 1 {
		t.Fatalf("unexpected metadata: %+v", got)
	}
	if len(got.Repositories) != 5 || len(got.Tools) != 3 {
		t.Fatalf("unexpected inventory: %+v", got)
	}
	if resourceCalls != 2 || containerCalls != 1 || !got.CollectedAt.Equal(wantTime) || got.Container.Name != ContainerName {
		t.Fatalf("unexpected fake observations: resources=%d container=%d snapshot=%+v", resourceCalls, containerCalls, got)
	}
	for _, v := range got.Repositories {
		if v.Present {
			t.Errorf("unexpected resource: %s", v.ID)
		}
	}
}

func TestDetectExecutable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test-tool")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0700); err != nil {
		t.Fatal(err)
	}
	got := inspectResources([]Target{{"a", "Tool", path, "executable"}})
	if len(got) != 1 || !got[0].Present {
		t.Fatalf("expected executable: %+v", got)
	}
}

func TestGuardBlocksHostAndMutations(t *testing.T) {
	h := LoopbackHostGuard(API(), "18484")
	cases := []struct {
		host, method string
		want         int
	}{
		{"127.0.0.1:18484", "GET", 200},
		{"localhost:18484", "GET", 200},
		{"attacker.example:18484", "GET", 403},
		{"127.0.0.1:18484", "POST", 403},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, "http://"+tc.host+"/api/v1/health", nil)
		r := httptest.NewRecorder()
		h.ServeHTTP(r, req)
		if r.Code != tc.want {
			t.Errorf("%s %s got %d want %d", tc.method, tc.host, r.Code, tc.want)
		}
	}
}

func TestGuardBlocksCrossOrigin(t *testing.T) {
	h := LoopbackHostGuard(API(), "18484")
	req := httptest.NewRequest("GET", "http://127.0.0.1:18484/api/v1/health", nil)
	req.Header.Set("Origin", "http://attacker.example")
	r := httptest.NewRecorder()
	h.ServeHTTP(r, req)
	if r.Code != 403 {
		t.Fatalf("got %d", r.Code)
	}
}
