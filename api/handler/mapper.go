package handler

import (
	"errors"
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/domain"
)

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error  string       `json:"error"`
	Detail string       `json:"detail,omitempty"`
	Fields []FieldError `json:"fields,omitempty"`
}

func mapCreateProjectError(err error) (int, ErrorResponse) {
	fields := createProjectFieldsValidation(err)
	if len(fields) > 0 {
		return http.StatusBadRequest, ErrorResponse{
			Error:  "validation_error",
			Detail: "validation failed",
			Fields: fields,
		}
	}

	if errors.Is(err, domain.ErrProjectExists) {
		return http.StatusConflict, ErrorResponse{
			Error:  "project_exists",
			Detail: domain.ErrProjectExists.Error(),
		}
	}

	return http.StatusInternalServerError, ErrorResponse{
		Error:  "internal_error",
		Detail: domain.ErrInternal.Error(),
	}
}

func mapGetProjectByIDError(err error) (int, ErrorResponse) {
	if errors.Is(err, domain.ErrProjectNotFound) {
		return http.StatusNotFound, ErrorResponse{
			Error:  "project_not_found",
			Detail: domain.ErrProjectNotFound.Error(),
		}
	}

	return http.StatusInternalServerError, ErrorResponse{
		Error:  "internal_error",
		Detail: domain.ErrInternal.Error(),
	}
}

var createProjectFieldsErrors = []struct {
	err   error
	field string
	code  string
}{
	{err: domain.ErrUnknownCategory, field: "category", code: "invalid_category"},
	{err: domain.ErrInvalidNameLen, field: "name", code: "invalid_project_name_length"},
	{err: domain.ErrInvalidDescriptionLen, field: "description", code: "invalid_project_description_length"},
	{err: domain.ErrUnknownCurrency, field: "currency", code: "invalid_currency"},
	{err: domain.ErrInvalidGoalAmount, field: "goalAmount", code: "invalid_goal_amount"},
	{err: domain.ErrInvalidDurationDays, field: "durationDays", code: "invalid_duration_days"},
}

func createProjectFieldsValidation(err error) []FieldError {
	fields := make([]FieldError, 0, len(createProjectFieldsErrors))
	for _, item := range createProjectFieldsErrors {
		if errors.Is(err, item.err) {
			fields = append(fields, FieldError{
				Field:   item.field,
				Code:    item.code,
				Message: item.err.Error(),
			})
		}
	}
	return fields
}
