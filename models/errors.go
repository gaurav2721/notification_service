package models

import "errors"

// Template-related errors
var (
	ErrInvalidTemplateContent  = errors.New("invalid template content")
	ErrInvalidTemplateType     = errors.New("invalid template type")
	ErrMissingRequiredVariable = errors.New("missing required variable")
	ErrTemplateNotFound        = errors.New("template not found")
)
