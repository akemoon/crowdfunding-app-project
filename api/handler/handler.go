package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/akemoon/golib/httplib"
	"github.com/google/uuid"
)

// CreateProject godoc
// @Summary      Create project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateProjectReq  true  "Project data"
// @Success      201   {object}  domain.Project
// @Failure      400   {object}  object
// @Failure      409   {object}  object
// @Failure      500   {object}  object
// @Router       /projects [post]
func CreateProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.CreateProjectReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		p, err := svc.CreateProject(r.Context(), req)
		if err != nil {
			log.Printf("CreateProject: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, CreateProjectMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusCreated, p)
	}
}

// GetProjectByID godoc
// @Summary      Get project by ID
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project UUID"
// @Success      200  {object}  domain.Project
// @Failure      400  {object}  object
// @Failure      404  {object}  object
// @Router       /projects/{id} [get]
func GetProjectByID(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		p, err := svc.GetProjectByID(r.Context(), id)
		if err != nil {
			log.Printf("GetProjectByID: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, GetProjectByIDMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, p)
	}
}

// GetProjects godoc
// @Summary      List all projects
// @Tags         projects
// @Produce      json
// @Success      200  {array}   domain.Project
// @Failure      500  {object}  object
// @Router       /projects [get]
func GetProjects(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := svc.GetProjects(r.Context())
		if err != nil {
			log.Printf("GetProjects: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, projects)
	}
}
