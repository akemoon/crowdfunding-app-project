package httplib

import (
	"errors"
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/golib/validation"
)

// TODO: move to golib

type ErrResp struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type ErrMapRule struct {
	Err     error
	Status  int
	Code    string
	Message string
}

// MapErrToHTTP maps an error to an HTTP status code and response body.
// *validation.Error is always handled automatically: returns 400 with field-level details.
// Other errors are matched against the provided rules in order.
func MapErrToHTTP(err error, rules []ErrMapRule) (int, ErrResp) {
	var ve *validation.Error
	if errors.As(err, &ve) {
		return http.StatusBadRequest, ErrResp{
			Code:    "validation_error",
			Message: "Validation failed",
			Fields:  ve.Fields(),
		}
	}

	for _, r := range rules {
		if errors.Is(err, r.Err) {
			return r.Status, ErrResp{Code: r.Code, Message: r.Message}
		}
	}

	return http.StatusInternalServerError, ErrResp{
		Code:    "internal_error",
		Message: "Unknown error",
	}
}