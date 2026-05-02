package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/golib/httplib"
	"github.com/akemoon/crowdfunding-app-project/service/project"
	"github.com/google/uuid"
)

// TODO: extract auth middleware for X-User-ID and X-User-Role checks
// to avoid duplicating the auth/role guard logic in every handler.

const (
	userIDHeader   = "X-User-ID"
	userRoleHeader = "X-User-Role"
	moderatorRole  = "moder"
)

type RejectProjectReq struct {
	Reason string `json:"reason"`
}

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

func GetProjectsByUserID(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		projects, err := svc.GetProjectsByUserID(r.Context(), userID, false)
		if err != nil {
			log.Printf("GetProjectsByUserID: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, projects)
	}
}

func GetMyProjects(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		projects, err := svc.GetProjectsByUserID(r.Context(), userID, true)
		if err != nil {
			log.Printf("GetMyProjects: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, projects)
	}
}

func GetProjectByID(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		var callerID *uuid.UUID
		parsedID, parseErr := httplib.ParseUUIDHeader(r, userIDHeader)
		if parseErr == nil {
			callerID = &parsedID
		}

		p, err := svc.GetProjectByID(r.Context(), id, callerID)
		if err != nil {
			log.Printf("GetProjectByID: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, GetProjectByIDMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, p)
	}
}

func ApproveProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != moderatorRole {
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
			status, errResp := httplib.MapErrToHTTP(err, ApproveProjectMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func RejectProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != moderatorRole {
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
			status, errResp := httplib.MapErrToHTTP(err, RejectProjectMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func GetMyApplications(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		managerID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != moderatorRole {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		apps, err := svc.GetMyApplications(r.Context(), managerID)
		if err != nil {
			log.Printf("GetMyApplications: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, apps)
	}
}

func TakeApplication(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		managerID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != moderatorRole {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		err = svc.TakeApplication(r.Context(), id, managerID)
		if err != nil {
			log.Printf("TakeApplication: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, TakeApplicationMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetApplicationByProjectID(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		app, err := svc.GetApplicationByProjectID(r.Context(), id, userID)
		if err != nil {
			log.Printf("GetApplicationByProjectID: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, GetApplicationByProjectIDMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, app)
	}
}

func GetPendingApplications(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Header.Get(userRoleHeader) != moderatorRole {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		apps, err := svc.GetPendingApplications(r.Context())
		if err != nil {
			log.Printf("GetPendingApplications: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, nil)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		httplib.WriteJSON(w, http.StatusOK, apps)
	}
}

type BoostProjectReq struct {
	PromoCode string `json:"promoCode"`
}

func BoostProject(svc *project.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := httplib.ParseUUIDHeader(r, userIDHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		var req BoostProjectReq

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.BoostProject(r.Context(), userID, id, req.PromoCode)
		if err != nil {
			log.Printf("BoostProject: %v", err)
			status, errResp := httplib.MapErrToHTTP(err, BoostProjectMapRules)
			httplib.WriteJSON(w, status, errResp)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
