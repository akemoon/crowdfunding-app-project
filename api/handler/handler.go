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
