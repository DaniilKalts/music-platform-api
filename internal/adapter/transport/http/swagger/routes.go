package swagger

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// Раздача статики Swagger UI
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/swagger/index.html")
	})

	// Раздача файлов спецификации OpenAPI
	r.Route("/swagger", func(r chi.Router) {
		fs := http.StripPrefix("/swagger/", http.FileServer(http.Dir("api/v1")))
		r.Handle("/*", fs)
	})
}
