package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// NewErrorResponse creates a new standardized error response
func NewErrorResponse(status int, errorType, title, detail, instance string) *ErrorResponse {
	return &ErrorResponse{
		Type:     errorType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

// NewBadRequestError creates a 400 Bad Request error
func NewBadRequestError(detail, instance string) *ErrorResponse {
	return NewErrorResponse(
		http.StatusBadRequest,
		TypeBadRequest,
		TitleBadRequest,
		detail,
		instance,
	)
}

// NewValidationError creates a 400 Bad Request error with validation details
func NewValidationError(err error, instance string) *ErrorResponse {
	errResp := NewErrorResponse(
		http.StatusBadRequest,
		TypeValidationFailed,
		TitleValidationFailed,
		"The request contains invalid or missing fields",
		instance,
	)

	// Parse validator.ValidationErrors if applicable
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errResp.Errors = parseValidationErrors(validationErrors)
	} else {
		// Fallback for other validation errors
		errResp.Detail = err.Error()
	}

	return errResp
}

// NewNotFoundError creates a 404 Not Found error
func NewNotFoundError(resourceType, instance string) *ErrorResponse {
	detail := fmt.Sprintf("The requested %s was not found", resourceType)
	return NewErrorResponse(
		http.StatusNotFound,
		TypeNotFound,
		TitleNotFound,
		detail,
		instance,
	)
}

// NewConflictError creates a 409 Conflict error
func NewConflictError(detail, instance string) *ErrorResponse {
	return NewErrorResponse(
		http.StatusConflict,
		TypeConflict,
		TitleConflict,
		detail,
		instance,
	)
}

// NewUnprocessableEntityError creates a 422 Unprocessable Entity error
func NewUnprocessableEntityError(detail, instance string) *ErrorResponse {
	return NewErrorResponse(
		http.StatusUnprocessableEntity,
		TypeUnprocessableEntity,
		TitleUnprocessableEntity,
		detail,
		instance,
	)
}

// NewInternalError creates a 500 Internal Server Error
func NewInternalError(instance string) *ErrorResponse {
	return NewErrorResponse(
		http.StatusInternalServerError,
		TypeInternalError,
		TitleInternalError,
		"An unexpected error occurred. Please try again later.",
		instance,
	)
}

// NewDatabaseError creates an appropriate error based on the database error type
func NewDatabaseError(err error, instance string) *ErrorResponse {
	if err == nil {
		return nil
	}

	// Check for record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewNotFoundError("resource", instance)
	}

	// Check for duplicate key or unique constraint violations
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "UNIQUE constraint failed") ||
		strings.Contains(errMsg, "Duplicate entry") {
		return NewConflictError("A resource with the same unique identifier already exists", instance)
	}

	// Check for foreign key constraint violations
	if strings.Contains(errMsg, "foreign key constraint") ||
		strings.Contains(errMsg, "FOREIGN KEY constraint failed") {
		return NewUnprocessableEntityError("The request references a resource that does not exist", instance)
	}

	// Default to internal server error for unknown database errors
	return NewInternalError(instance)
}

// WriteError writes a JSON error response with proper headers
func WriteError(w http.ResponseWriter, errResp *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(errResp.Status)

	// Ignore encoding errors - we've already sent the status code
	_ = json.NewEncoder(w).Encode(errResp)
}

// parseValidationErrors converts validator.ValidationErrors to our ValidationError format
func parseValidationErrors(validationErrs validator.ValidationErrors) []ValidationError {
	errors := make([]ValidationError, 0, len(validationErrs))

	for _, err := range validationErrs {
		field := err.Field()
		tag := err.Tag()

		// Convert first character to lowercase for JSON field names
		if len(field) > 0 {
			field = strings.ToLower(field[:1]) + field[1:]
		}

		message := formatValidationMessage(field, tag, err.Param())

		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
			Code:    tag,
			Value:   err.Value(),
		})
	}

	return errors
}

// formatValidationMessage creates a human-readable validation error message
func formatValidationMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s", field, param)
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", field)
	case "url":
		return fmt.Sprintf("Field '%s' must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("Field '%s' must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("Field '%s' must be one of: %s", field, param)
	case "gt":
		return fmt.Sprintf("Field '%s' must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("Field '%s' must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("Field '%s' must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("Field '%s' must be less than or equal to %s", field, param)
	default:
		return fmt.Sprintf("Field '%s' failed validation on '%s' tag", field, tag)
	}
}
