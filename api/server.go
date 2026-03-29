package api

import (
	"net/http"

	_ "github.com/akemoon/crowdfunding-app-project/docs"
	"github.com/akemoon/crowdfunding-app-project/api/handler"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/akemoon/golib/httplib"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	s *http.Server
	r *httplib.Router
}

func NewServer() *Server {
	return &Server{
		r: httplib.NewRouter(),
	}
}

func (s *Server) AddProjectHandlers(svc *project.Service) {
	s.r.HandleFunc("POST /projects", handler.CreateProject(svc))
	s.r.HandleFunc("GET /projects", handler.GetProjects(svc))
	s.r.HandleFunc("GET /projects/{id}", handler.GetProjectByID(svc))
	s.r.Handle("/swagger/", httpSwagger.WrapHandler)
}

func (s *Server) ListenAndServe(addr string) error {
	s.s = &http.Server{
		Addr:    addr,
		Handler: s.r.Handler(),
	}
	return s.s.ListenAndServe()
}
