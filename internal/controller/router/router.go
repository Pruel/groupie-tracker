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

	// Можно сложить логику различных контроллеров в один и это будет более лучшая практика.
	r.Mux.HandleFunc("/", controller.MainController)         // main {filters, search} <- want do filter ->
	r.Mux.HandleFunc("/filter", controller.FilterController) // main {filters, s} <- also this feature
	r.Mux.HandleFunc("/search", controller.SearchController) // main   <- feature of the main page
	r.Mux.HandleFunc("/artist", controller.ArtistController)

	// Method = HTTP "/main" {GET, post, delete, put} = one controller per method (http) -> RESTful API + FRONT-END + BACK-END = two different service

	// Frontend_backend = programm, single service
	// one controller per tempalate

	slog.Info("the ServeMux router successful initialized")

	return nil
}
