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
	CreateProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectExists,
	}
	GetProjectByIDMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
	}
)
