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
	r.Mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("web/screenshot"))))

	// r.Mux.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir(controller.GetTmplFilepath("view")))))
	r.Mux.HandleFunc("/", controller.FilterController)       //
	r.Mux.HandleFunc("/filter", controller.FilterController) //
	r.Mux.HandleFunc("/artist", controller.ArtistController) //
	r.Mux.HandleFunc("/search", controller.SearchController)

	slog.Debug("the ServeMux router successful initialized")

	return nil
}
