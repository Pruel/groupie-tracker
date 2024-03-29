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

	r.Mux.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("./../../view/"))))
	http.Handle("./view/static", http.StripPrefix("./view/static/", http.FileServer(http.Dir("view/static"))))
	// r.Mux.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir(controller.GetTmplFilepath("view")))))
	r.Mux.HandleFunc("/", controller.MainController)         //
	r.Mux.HandleFunc("/artist", controller.ArtistController) //

	slog.Debug("the ServeMux router successful initialized")

	return nil
}
