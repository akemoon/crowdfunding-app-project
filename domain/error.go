package domain

import "errors"

var (
	ErrUnknownCategory = errors.New("unknown category")

	ErrInvalidNameLen = errors.New("invalid project name length")

	ErrInvalidDescriptionLen = errors.New("invalid project description length")

	ErrInvalidGoalAmount = errors.New("invalid goal amount")

	ErrInvalidDurationDays = errors.New("invalid duration days")

	ErrUnknownCurrency = errors.New("unknown currency")

	ErrInvalidDonationAmount = errors.New("invalid donation amount")

	ErrProjectNotFound = errors.New("project not found")

	ErrProjectExists = errors.New("project name already exists")

	ErrUnknownConflict = errors.New("unknown conflict")

	ErrInternal = errors.New("internal error")

	ErrUnknownStatus = errors.New("unknown status")

	ErrUnknownSort = errors.New("unknown sort")

	ErrForbidden = errors.New("forbidden")

	ErrProjectNotDraft = errors.New("project is not a draft")

	ErrProjectNotOnReview = errors.New("project is not on review")

	ErrApplicationNotFound = errors.New("application not found")

	ErrApplicationAlreadyTaken = errors.New("application is already taken")

	ErrProjectNotActive = errors.New("project is not active")

	ErrProjectImageNotFound = errors.New("project image not found")

	ErrUnsupportedFileType = errors.New("unsupported file type")

	ErrProjectCoverRequired = errors.New("project cover is required")
)
