package handler

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/client/promocode"
	"github.com/akemoon/crowdfunding-app-project/domain"
	"github.com/akemoon/golib/httplib"
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
	MapRuleProjectNotOnReview = httplib.ErrMapRule{
		Err:     domain.ErrProjectNotOnReview,
		Status:  http.StatusConflict,
		Code:    "project_not_on_review",
		Message: domain.ErrProjectNotOnReview.Error(),
	}
	MapRuleApplicationNotFound = httplib.ErrMapRule{
		Err:     domain.ErrApplicationNotFound,
		Status:  http.StatusNotFound,
		Code:    "application_not_found",
		Message: domain.ErrApplicationNotFound.Error(),
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
	MapRuleProjectImageNotFound = httplib.ErrMapRule{
		Err:     domain.ErrProjectImageNotFound,
		Status:  http.StatusNotFound,
		Code:    "project_image_not_found",
		Message: domain.ErrProjectImageNotFound.Error(),
	}
	MapRuleUnsupportedFileType = httplib.ErrMapRule{
		Err:     domain.ErrUnsupportedFileType,
		Status:  http.StatusBadRequest,
		Code:    "unsupported_file_type",
		Message: domain.ErrUnsupportedFileType.Error(),
	}
)

var (
	MapRulePromoCodeNotFound = httplib.ErrMapRule{
		Err:     promocode.ErrPromoCodeNotFound,
		Status:  http.StatusNotFound,
		Code:    "promo_code_not_found",
		Message: promocode.ErrPromoCodeNotFound.Error(),
	}
	MapRulePromoCodeUsed = httplib.ErrMapRule{
		Err:     promocode.ErrPromoCodeUsed,
		Status:  http.StatusConflict,
		Code:    "promo_code_already_used",
		Message: promocode.ErrPromoCodeUsed.Error(),
	}
	MapRulePromoCodeAccessDenied = httplib.ErrMapRule{
		Err:     promocode.ErrPromoCodeForbidden,
		Status:  http.StatusForbidden,
		Code:    "promo_code_access_denied",
		Message: promocode.ErrPromoCodeForbidden.Error(),
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
	ApproveProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotOnReview,
	}
	RejectProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotOnReview,
	}
	TakeApplicationMapRules = []httplib.ErrMapRule{
		MapRuleApplicationNotFound,
	}
	GetApplicationByProjectIDMapRules = []httplib.ErrMapRule{
		MapRuleApplicationNotFound,
	}
	UpdateProjectMapRules = []httplib.ErrMapRule{
		MapRuleApplicationNotFound,
		MapRuleProjectExists,
	}
	BoostProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRulePromoCodeNotFound,
		MapRulePromoCodeUsed,
		MapRulePromoCodeAccessDenied,
	}
	UploadProjectImageMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleUnsupportedFileType,
	}
	DeleteProjectImageMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleProjectImageNotFound,
	}
)
