package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/golib/httplib"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/google/uuid"
)

const userIDHeader = "X-User-ID"
const userRoleHeader = "X-User-Role"
const managerRole = "Manager"

type RejectProjectReq struct {
	Reason string `json:"reason"`
}

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
// @Router /projects [post]
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
			log.Printf("CreateProject: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, CreateProjectMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary Get projects
// @Description Returns a paginated list of projects filtered by status, category and search.
// @Tags project
// @Produce json
// @Param status query string false "Status filter (active|finished), default: active"
// @Param sort query string false "Sort order (default|date), default: default"
// @Param category query string false "Category filter (science|tech|architecture_and_urban|sport|music)"
// @Param search query string false "Search by name"
// @Param limit query int false "Page size (1-100, default: 20)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} domain.GetProjectsResp
// @Failure 400 {object} httplib.ErrResp "unknown_status | unknown_sort | unknown_category"
// @Failure 500 {object} httplib.ErrResp "internal_error"
// @Router /projects [get]
func GetProjects(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		status := domain.Status(q.Get("status"))
		if status == "" {
			status = domain.StatusActive
		}

		var category *domain.Category
		if c := q.Get("category"); c != "" {
			cat := domain.Category(c)
			category = &cat
		}

		var search *string
		if s := q.Get("search"); s != "" {
			search = &s
		}

		sort := domain.Sort(q.Get("sort"))

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		resp, err := svc.GetProjects(r.Context(), domain.GetProjectsReq{
			Status:   status,
			Sort:     sort,
			Category: category,
			Search:   search,
			Limit:    limit,
			Offset:   offset,
		})
		if err != nil {
			log.Printf("GetProjects: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, GetProjectsMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, resp)
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
// @Router /projects/{id} [get]
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

// @Summary Approve project application
// @Description Approves a project application and moves project to active.
// @Tags manager
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param X-User-Role header string true "User role (must be Manager)"
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]string "status=ok"
// @Failure 400 {string} string "invalid project id"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 409 {object} httplib.ErrResp "project_not_on_review"
// @Failure 500 {object} httplib.ErrResp "internal_error"
// @Router /manager/applications/{id}/approve [post]
func ApproveProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != managerRole {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		err = svc.ApproveProject(r.Context(), id)
		if err != nil {
			log.Printf("ApproveProject: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// @Summary Reject project application
// @Description Rejects a project application with reason.
// @Tags manager
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param X-User-Role header string true "User role (must be Manager)"
// @Param id path string true "Project ID"
// @Param body body RejectProjectReq true "Reject reason payload"
// @Success 200 {object} map[string]string "status=ok"
// @Failure 400 {string} string "invalid project id | invalid request body"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 409 {object} httplib.ErrResp "project_not_on_review"
// @Failure 500 {object} httplib.ErrResp "internal_error"
// @Router /manager/applications/{id}/reject [post]
func RejectProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != managerRole {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		var req RejectProjectReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.RejectProject(r.Context(), id, req.Reason)
		if err != nil {
			log.Printf("RejectProject: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// TODO: boost project req -> sync to promo
