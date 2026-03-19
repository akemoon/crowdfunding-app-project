package api

import (
	"net/http"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/akemoon/crowdfunding-app-project/api/handler"
	"github.com/akemoon/crowdfunding-app-project/graph"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/akemoon/golib/httplib"
	"github.com/akemoon/golib/httplib/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	s *http.Server
	r *httplib.Router
}

func NewServer() *Server {
	return &Server{
		r: httplib.NewRouter().Use(
			middleware.BaseMetrics(),
		),
	}
}

func (s *Server) AddProjectHandlers(svc *project.Service) {
	s.r.HandleFunc("POST /projects", handler.CreateProject(svc))
	s.r.HandleFunc("GET /projects", handler.GetProjects(svc))
	s.r.HandleFunc("GET /projects/user", handler.GetMyProjects(svc))
	s.r.HandleFunc("GET /projects/user/{id}", handler.GetProjectsByUserID(svc))
	s.r.HandleFunc("GET /projects/{id}", handler.GetProjectByID(svc))
	s.r.HandleFunc("POST /projects/{id}/boost", handler.BoostProject(svc))

	// TODO: move application handlers to a separate registry (applications are a distinct resource)
	s.r.HandleFunc("GET /applications", handler.GetPendingApplications(svc))
	s.r.HandleFunc("GET /applications/moderator", handler.GetMyApplications(svc))
	s.r.HandleFunc("GET /applications/project/{id}", handler.GetApplicationByProjectID(svc))
	s.r.HandleFunc("POST /applications/{id}/take", handler.TakeApplication(svc))
	s.r.HandleFunc("POST /applications/{id}/approve", handler.ApproveProject(svc))
	s.r.HandleFunc("POST /applications/{id}/reject", handler.RejectProject(svc))
}

func (s *Server) AddGraphQL(svc *project.Service) {
	srv := gqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(svc),
	}))
	s.r.Handle("/graphql", srv)
	s.r.Handle("/playground", playground.Handler("GraphQL", "/graphql"))
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
