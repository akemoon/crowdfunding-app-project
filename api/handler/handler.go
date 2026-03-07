package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/google/uuid"
)

// @Summary Create project
// @Description Creates a new project by given payload.
// @Tags project
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param body body domain.CreateProjectReq true "Project create payload"
// @Success 201 {string} string "created"
// @Failure 400 {object} ErrorResponse "validation_error | invalid request body"
// @Failure 401 {string} string "unauthorized"
// @Failure 405 {string} string "method not allowed"
// @Failure 409 {object} ErrorResponse "project_exists"
// @Failure 500 {object} ErrorResponse "internal_error"
// @Router /project [post]
func CreateProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
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
			status, resp := mapCreateProjectError(err)
			writeJSON(w, status, resp)
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
// @Failure 404 {object} ErrorResponse "project_not_found"
// @Failure 405 {string} string "method not allowed"
// @Failure 500 {object} ErrorResponse "internal_error"
// @Router /project/{id} [get]
func GetProjectByID(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.PathValue("id")

		// TODO: check empty string
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		p, err := svc.GetProjectByID(r.Context(), id)
		if err != nil {
			status, resp := mapGetProjectByIDError(err)
			writeJSON(w, status, resp)
			return
		}

		writeJSON(w, http.StatusOK, p)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
