package api

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/api/handler"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/akemoon/golib/myhttp"
	"github.com/akemoon/golib/myhttp/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	s *http.Server
	r *myhttp.Router
}

func NewServer() *Server {
	return &Server{
		r: myhttp.NewRouter().Use(
			middleware.BaseMetrics(),
		),
	}
}

func (s *Server) AddProjectHandlers(svc *project.Service) {
	s.r.HandleFunc("POST /projects", handler.CreateProject(svc))
	s.r.HandleFunc("GET /projects", handler.GetProjects(svc))
	s.r.HandleFunc("GET /projects/{id}", handler.GetProjectByID(svc))
	s.r.HandleFunc("GET /projects/{id}/applications", handler.GetApplicationByProjectID(svc))
	s.r.HandleFunc("POST /projects/{id}/boost", handler.BoostProject(svc))

	s.r.HandleFunc("GET /projects/applications", handler.GetPendingApplications(svc))
	s.r.HandleFunc("GET /projects/applications/moderator", handler.GetMyApplications(svc))
	s.r.HandleFunc("POST /projects/applications/{id}/take", handler.TakeApplication(svc))
	s.r.HandleFunc("POST /projects/applications/{id}/approve", handler.ApproveProject(svc))
	s.r.HandleFunc("POST /projects/applications/{id}/reject", handler.RejectProject(svc))
}

func (s *Server) AddSwaggerUI() {
	s.r.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

func (s *Server) AddMetrics() {
	s.r.Handle("/metrics", promhttp.Handler())
}

func (s *Server) ListenAndServe(addr string) error {
	s.s = &http.Server{
		Addr:    addr,
		Handler: s.r.Handler(),
	}
	return s.s.ListenAndServe()
}

//func (s *Server) Stop(ctx context.Context) error {
//	s.s.Shutdown()
//}
