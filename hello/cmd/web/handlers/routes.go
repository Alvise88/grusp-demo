package handlers

import (
	"context"
	"net/http"

	"schneider.vip/problem"

	"github.com/alexedwards/flow"
	"grusp.io/hello/cmd/web/json"
	"grusp.io/hello/pkg/model"
)

func Routes(ctx context.Context) *flow.Mux {
	// Initialize a new router.
	mux := flow.New()

	// pipelines
	mux.HandleFunc("/health", Health(ctx), http.MethodGet)

	return mux
}

func Health(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		status := model.Health{
			Status: "ok",
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.Write(w, http.StatusOK, json.Envelope{"status": status}, nil)

		if err != nil {
			p := problem.New(
				problem.Type("https://argoci/500"),
				problem.Status(http.StatusInternalServerError),
			)
			_ = json.Write(w, http.StatusInternalServerError, json.Envelope{"problem": p}, nil)
			return
		}
	}
}
