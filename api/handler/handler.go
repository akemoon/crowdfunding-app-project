package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/golib/httplib"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/google/uuid"
)

const userIDHeader = "X-User-ID"

// @Summary Create project
// @Description Creates a new project by given payload.
// @Tags project
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param body body domain.CreateProjectReq true "Project create payload"
// @Success 201 {string} string "created"
// @Failure 400 {object} httplib.ErrResp "validation_error | invalid request body"
// @Failure 401 {string} string "unauthorized"
// @Failure 405 {string} string "method not allowed"
// @Failure 409 {object} httplib.ErrResp "project_exists"
// @Failure 500 {object} httplib.ErrResp "internal_error"
// @Router /project [post]
func CreateProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req domain.CreateProjectReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.CreateProject(r.Context(), userID, req)
		if err != nil {
			status, errResp := httplib.MapErrToHTTP(err, CreateProjectMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary Get project by id
// @Description Returns a project by its id.
// @Tags project
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Project
// @Failure 400 {string} string "invalid project id"
// @Failure 404 {object} httplib.ErrResp "project_not_found"
// @Failure 405 {string} string "method not allowed"
// @Failure 500 {object} httplib.ErrResp "internal_error"
// @Router /project/{id} [get]
func GetProjectByID(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		p, err := svc.GetProjectByID(r.Context(), id)
		if err != nil {
			status, errResp := httplib.MapErrToHTTP(err, GetProjectByIDMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, p)
	}
}
