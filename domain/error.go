package domain

import "errors"

var (
	ErrUnknownCategory = errors.New("unknown category")

	ErrInvalidNameLen = errors.New("invalid project name length")

	ErrInvalidDescriptionLen = errors.New("invalid project description length")

	ErrInvalidGoalAmount = errors.New("invalid goal amount")

	ErrInvalidDurationDays = errors.New("invalid duration days")

	ErrUnknownCurrency = errors.New("unknown currency")

	ErrProjectNotFound = errors.New("project not found")

	ErrProjectExists = errors.New("project name already exists")

	ErrInternal = errors.New("internal error")
)
