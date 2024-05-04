package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/phillipmugisa/sermon_finder/monitor"
)

type Server struct {
	port      string
	storage   *sql.DB
	logger    *slog.Logger
	templates *template.Template
	monitor   *monitor.RequestMonitor
}

func NewServer(port string, storage *sql.DB, logger *slog.Logger) *Server {
	templates := NewTemplate()
	monitor := monitor.NewRequestMethod()

	return &Server{
		port:      port,
		storage:   storage,
		logger:    logger,
		templates: templates,
		monitor:   monitor,
	}
}

func (server *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", server.port),
		Handler:        r,
		ReadTimeout:    10 & time.Second,
		WriteTimeout:   10 & time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.serverStatic(r)
	server.registerRoutes(r)

	server.logger.Info(fmt.Sprintf("starting server on port %s ...\n", server.port))
	return s.ListenAndServe()
}

func (server *Server) registerRoutes(r *chi.Mux) {
	r.Get("/", server.MakeRequestHandler(server.HomeHandler))
	r.Post("/sermon/upload/", server.MakeRequestHandler(server.SermonUploadHandler))

}

func (server *Server) serverStatic(r *chi.Mux) {
	staticfileserver := http.FileServer(http.Dir("./static/"))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", staticfileserver)
		fs.ServeHTTP(w, r)
	})
}
