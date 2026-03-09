package handler

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/crowdfunding-app-project/golib/httplib"
)

var (
	MapRuleProjectExists = httplib.ErrMapRule{
		Err:     domain.ErrProjectExists,
		Status:  http.StatusConflict,
		Code:    "project_exists",
		Message: domain.ErrProjectExists.Error(),
	}
	MapRuleProjectNotFound = httplib.ErrMapRule{
		Err:     domain.ErrProjectNotFound,
		Status:  http.StatusNotFound,
		Code:    "project_not_found",
		Message: domain.ErrProjectNotFound.Error(),
	}
)

var (
	MapRuleUnknownStatus = httplib.ErrMapRule{
		Err:     domain.ErrUnknownStatus,
		Status:  http.StatusBadRequest,
		Code:    "unknown_status",
		Message: domain.ErrUnknownStatus.Error(),
	}
	MapRuleUnknownCategory = httplib.ErrMapRule{
		Err:     domain.ErrUnknownCategory,
		Status:  http.StatusBadRequest,
		Code:    "unknown_category",
		Message: domain.ErrUnknownCategory.Error(),
	}
	MapRuleUnknownSort = httplib.ErrMapRule{
		Err:     domain.ErrUnknownSort,
		Status:  http.StatusBadRequest,
		Code:    "unknown_sort",
		Message: domain.ErrUnknownSort.Error(),
	}
)

var (
	CreateProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectExists,
	}
	GetProjectByIDMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
	}
	GetProjectsMapRules = []httplib.ErrMapRule{
		MapRuleUnknownStatus,
		MapRuleUnknownCategory,
		MapRuleUnknownSort,
	}
)
