package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/kol-workbench/internal/host"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

const port = "18484"

type dashboard struct {
	app.Compo
	snapshot host.Snapshot
	loaded   bool
	loading  bool
	failure  string
}

func (d *dashboard) OnMount(ctx app.Context) { d.refresh(ctx) }

func (d *dashboard) refresh(ctx app.Context) {
	d.loading = true
	d.failure = ""
	go func() {
		client := http.Client{Timeout: 7 * time.Second}
		resp, err := client.Get("/api/v1/status")
		var snapshot host.Snapshot
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = fmt.Errorf("gateway HTTP %d", resp.StatusCode)
			} else {
				err = json.NewDecoder(resp.Body).Decode(&snapshot)
			}
		}
		ctx.Dispatch(func(ctx app.Context) {
			d.loading = false
			if err != nil {
				d.failure = err.Error()
				return
			}
			d.loaded = true
			d.snapshot = snapshot
		})
	}()
}

func stateLabel(present bool) string {
	if present {
		return "present"
	}
	return "missing"
}

func (d *dashboard) Render() app.UI {
	repoCards := make([]app.UI, 0, len(d.snapshot.Repositories))
	for _, r := range d.snapshot.Repositories {
		repoCards = append(repoCards, resourceCard(r))
	}
	toolCards := make([]app.UI, 0, len(d.snapshot.Tools))
	for _, t := range d.snapshot.Tools {
		toolCards = append(toolCards, resourceCard(t))
	}
	c := d.snapshot.Container
	status := "unknown"
	statusClass := "neutral"
	if c.Available {
		status = c.Status
		if c.Running {
			statusClass = "present"
		} else {
			statusClass = "stopped"
		}
	}

	return app.Div().Class("shell").Body(
		app.Header().Class("topbar").Body(
			app.Div().Body(
				app.P().Class("eyebrow").Text("LOCAL DEVELOPMENT · READ-ONLY"),
				app.H1().Text("KoL Workbench"),
				app.P().Class("subtitle").Text("Host resources, tools and the existing KoLmafia container — no mutations."),
			),
			app.Button().Class("refresh").Text("Refresh inventory").OnClick(func(ctx app.Context, e app.Event) { d.refresh(ctx) }),
		),
		app.If(d.loading, func() app.UI { return app.P().Class("notice").Text("Checking local resources...") }),
		app.If(d.failure != "", func() app.UI { return app.P().Class("error").Text("Gateway error: " + d.failure) }),
		app.If(!d.loaded, func() app.UI { return app.P().Class("notice").Text("Waiting for first inventory snapshot.") }),
		app.If(d.loaded, func() app.UI {
			return app.Main().Body(
				app.Section().Class("overview").Body(
					app.Div().Class("panel").Body(
						app.P().Class("eyebrow").Text("RUNTIME"),
						app.H2().Text(c.Name),
						app.P().Class("pill "+statusClass).Text(status),
						app.P().Class("muted").Text(fmt.Sprintf("Docker inspect only · running=%v · exit=%d", c.Running, c.ExitCode)),
					),
					app.Div().Class("panel").Body(
						app.P().Class("eyebrow").Text("GUARDRAILS"),
						app.H2().Text("Observe, do not execute"),
						app.P().Class("pill present").Text(d.snapshot.Mode),
						app.P().Class("muted").Text("No CLI execution, build jobs, volume changes or container lifecycle endpoints."),
					),
				),
				app.Section().Body(
					app.H2().Class("section-title").Text("Registered repositories / indexes"),
					app.Div().Class("cards").Body(repoCards...),
				),
				app.Section().Body(
					app.H2().Class("section-title").Text("Known local CLI tools"),
					app.Div().Class("cards").Body(toolCards...),
				),
				app.Footer().Body(
					app.P().Class("muted").Text("Inventory at "+d.snapshot.CollectedAt.Local().Format(time.RFC3339)),
					app.P().Class("muted").Text("Only these known targets are inspected. Missing does not mean deleted."),
				),
			)
		}),
	)
}

func resourceCard(r host.Resource) app.UI {
	s := stateLabel(r.Present)
	return app.Div().Class("panel resource").Body(
		app.P().Class("eyebrow").Text(r.Kind),
		app.H3().Text(r.Label),
		app.P().Class("pill "+s).Text(s),
		app.P().Class("path").Text(r.Path),
	)
}

func main() {
	app.Route("/", func() app.Composer { return &dashboard{} })
	app.RunWhenOnBrowser()

	mux := http.NewServeMux()
	mux.Handle("/api/", host.API())
	mux.Handle("/", &app.Handler{
		Name:        "KoL Workbench",
		Description: "Local KoLmafia development resource dashboard (read-only)",
		Styles:      []string{"/web/style.css"},
	})
	server := &http.Server{
		Addr:              "127.0.0.1:" + port,
		Handler:           host.LoopbackHostGuard(mux, port),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("read-only dashboard: http://127.0.0.1:%s", port)
	log.Fatal(server.ListenAndServe())
}
