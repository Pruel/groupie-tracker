package router

import (
	"log/slog"
	"net/http"

	"groupie-tracker/internal/controller"
)

// custom router
type Router struct {
	Mux *http.ServeMux // ServeMux
}

// New create and customazing the ServeMux router, and return
func New() *Router {
	return &Router{
		Mux: http.NewServeMux(),
	}
}

// InitRouter ...
func (r *Router) InitRouter() error {
	slog.Debug("the ServeMux router initializing")

	r.Mux.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("view"))))
	r.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/view/static"))))
	r.Mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("internal/view/static/image"))))

	r.Mux.HandleFunc("/", controller.MainController)
	r.Mux.HandleFunc("/filter", controller.FilterController)
	r.Mux.HandleFunc("/artist", controller.ArtistController)
	r.Mux.HandleFunc("/search", controller.SearchController)

	slog.Debug("the ServeMux router successful initialized")

	return nil
}
