package handlers

import (
	"context"
	"net/http"
	"os"

	"schneider.vip/problem"

	"github.com/alexedwards/flow"
	"grusp.io/hello/cmd/web/application"
	"grusp.io/hello/cmd/web/json"
	"grusp.io/hello/internal/models"
	"grusp.io/hello/pkg/model"
	"grusp.io/hello/ui"
)

func Routes(ctx context.Context, app *application.Application) *flow.Mux {
	// Initialize a new router.
	mux := flow.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("/static/...", fileServer)

	// healthcheck
	mux.HandleFunc("/health", health(ctx), http.MethodGet)

	// home
	mux.HandleFunc("/", home(app), http.MethodGet)

	return mux
}

func home(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := models.Info{
			Namespace: "default",
			Node:      "-",
			Pod:       "-",
			Container: "-",
			Message:   "Hello Grusp!",
		}

		namespace, ok := os.LookupEnv("KUBERNETES_NAMESPACE")

		if ok {
			info.Namespace = namespace
		}

		node, ok := os.LookupEnv("KUBERNETES_NODE_NAME")

		if ok {
			info.Node = node
		}

		pod, ok := os.LookupEnv("KUBERNETES_POD_NAME")

		if ok {
			info.Pod = pod
		}

		container, ok := os.LookupEnv("CONTAINER_IMAGE")

		if ok {
			info.Container = container
		}

		data := &models.TemplateData{
			CurrentYear: 2023,

			Info: info,
		}

		app.Render(w, http.StatusOK, "home.tmpl", data)
	}
}

func health(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := model.Health{
			Status: "ok",
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.Write(w, http.StatusOK, json.Envelope{"status": status}, nil)

		if err != nil {
			p := problem.New(
				problem.Type("https://grusp/500"),
				problem.Status(http.StatusInternalServerError),
			)
			_ = json.Write(w, http.StatusInternalServerError, json.Envelope{"problem": p}, nil)
			return
		}
	}
}
