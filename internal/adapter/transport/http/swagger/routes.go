package swagger

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	// Single wildcard handler so the UI shell and the spec files never collide:
	// "/swagger/" serves the Swagger UI, every other path is served from the
	// OpenAPI spec dir (openapi.yaml, paths/, components/).
	specFiles := http.StripPrefix("/swagger/", http.FileServer(http.Dir("api/v1")))
	r.Get("/swagger/*", func(w http.ResponseWriter, req *http.Request) {
		// The spec is multi-file and edited often; never let the browser serve a
		// stale openapi.yaml / path file, or $ref resolution breaks silently.
		w.Header().Set("Cache-Control", "no-store, must-revalidate")
		if req.URL.Path == "/swagger/" {
			http.ServeFile(w, req, "web/swagger/index.html")
			return
		}
		specFiles.ServeHTTP(w, req)
	})
}
