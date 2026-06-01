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
	MapRuleProjectNotDraft = httplib.ErrMapRule{
		Err:     domain.ErrProjectNotDraft,
		Status:  http.StatusConflict,
		Code:    "project_not_draft",
		Message: domain.ErrProjectNotDraft.Error(),
	}
	MapRuleProjectNotActive = httplib.ErrMapRule{
		Err:     domain.ErrProjectNotActive,
		Status:  http.StatusConflict,
		Code:    "project_not_active",
		Message: domain.ErrProjectNotActive.Error(),
	}
	MapRuleProjectImageNotFound = httplib.ErrMapRule{
		Err:     domain.ErrProjectImageNotFound,
		Status:  http.StatusNotFound,
		Code:    "project_image_not_found",
		Message: domain.ErrProjectImageNotFound.Error(),
	}
	MapRuleApplicationNotFound = httplib.ErrMapRule{
		Err:     domain.ErrApplicationNotFound,
		Status:  http.StatusNotFound,
		Code:    "application_not_found",
		Message: domain.ErrApplicationNotFound.Error(),
	}
	MapRuleForbidden = httplib.ErrMapRule{
		Err:     domain.ErrForbidden,
		Status:  http.StatusForbidden,
		Code:    "forbidden",
		Message: domain.ErrForbidden.Error(),
	}
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
	MapRuleUnsupportedFileType = httplib.ErrMapRule{
		Err:     domain.ErrUnsupportedFileType,
		Status:  http.StatusBadRequest,
		Code:    "unsupported_file_type",
		Message: domain.ErrUnsupportedFileType.Error(),
	}
	MapRuleFileTooLarge = httplib.ErrMapRule{
		Err:     domain.ErrFileTooLarge,
		Status:  http.StatusRequestEntityTooLarge,
		Code:    "file_too_large",
		Message: domain.ErrFileTooLarge.Error(),
	}
	MapRuleProjectCoverRequired = httplib.ErrMapRule{
		Err:     domain.ErrProjectCoverRequired,
		Status:  http.StatusUnprocessableEntity,
		Code:    "project_cover_required",
		Message: domain.ErrProjectCoverRequired.Error(),
	}
	MapRulePromoCodeNotFound = httplib.ErrMapRule{
		Err:     promocode.ErrPromoCodeNotFound,
		Status:  http.StatusNotFound,
		Code:    "promo_code_not_found",
		Message: promocode.ErrPromoCodeNotFound.Error(),
	}
	MapRuleInvalidCodeFormat = httplib.ErrMapRule{
		Err:     promocode.ErrInvalidCodeFormat,
		Status:  http.StatusBadRequest,
		Code:    "invalid_promo_code_format",
		Message: promocode.ErrInvalidCodeFormat.Error(),
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
		MapRuleProjectNotFound,
		MapRuleForbidden,
		MapRuleProjectNotDraft,
		MapRuleProjectExists,
	}
	BoostProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleForbidden,
		MapRuleProjectNotActive,
		MapRuleInvalidCodeFormat,
		MapRulePromoCodeNotFound,
	}
	SubmitProjectMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleForbidden,
		MapRuleProjectNotDraft,
		MapRuleProjectCoverRequired,
	}
	UploadProjectCoverMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleForbidden,
		MapRuleProjectNotDraft,
		MapRuleFileTooLarge,
		MapRuleUnsupportedFileType,
	}
	UploadProjectImageMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleForbidden,
		MapRuleProjectNotDraft,
		MapRuleFileTooLarge,
		MapRuleUnsupportedFileType,
	}
	DeleteProjectImageMapRules = []httplib.ErrMapRule{
		MapRuleProjectNotFound,
		MapRuleProjectImageNotFound,
	}
)
