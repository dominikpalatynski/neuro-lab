package errors

// ErrorResponse represents a standardized error response following RFC 7807 Problem Details
type ErrorResponse struct {
	// Type is a URI reference that identifies the problem type
	Type string `json:"type"`
	// Title is a short, human-readable summary of the problem type
	Title string `json:"title"`
	// Status is the HTTP status code
	Status int `json:"status"`
	// Detail is a human-readable explanation specific to this occurrence
	Detail string `json:"detail"`
	// Instance is a URI reference that identifies the specific occurrence
	Instance string `json:"instance"`
	// Errors contains validation errors if applicable
	Errors []ValidationError `json:"errors,omitempty"`
	// Meta contains additional metadata about the error
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// ValidationError represents a field-level validation error
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string `json:"field"`
	// Message is a human-readable error message
	Message string `json:"message"`
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Value is the rejected value (optional, for debugging)
	Value interface{} `json:"value,omitempty"`
}

// Error type constants following OpenAPI conventions
const (
	// TypeValidationFailed indicates request validation failure
	TypeValidationFailed = "validation-failed"
	// TypeNotFound indicates requested resource not found
	TypeNotFound = "not-found"
	// TypeConflict indicates resource conflict (duplicate, constraint violation)
	TypeConflict = "conflict"
	// TypeInternalError indicates internal server error
	TypeInternalError = "internal-error"
	// TypeBadRequest indicates malformed request
	TypeBadRequest = "bad-request"
	// TypeUnprocessableEntity indicates semantically invalid request
	TypeUnprocessableEntity = "unprocessable-entity"
)

// Error titles for consistent messaging
const (
	TitleValidationFailed    = "Validation Failed"
	TitleNotFound            = "Resource Not Found"
	TitleConflict            = "Resource Conflict"
	TitleInternalError       = "Internal Server Error"
	TitleBadRequest          = "Bad Request"
	TitleUnprocessableEntity = "Unprocessable Entity"
)
