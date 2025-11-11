package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"config/errors"

	"github.com/google/uuid"
)

// RequestIDKey is the context key for request IDs
type contextKey string

const RequestIDKey contextKey = "requestID"

// RequestID middleware adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID is already present in header
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate new UUID for request ID
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// Recovery middleware recovers from panics and returns a standardized error response
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Get request ID from context
				requestID := GetRequestID(r.Context())

				// Create error response
				errResp := errors.NewInternalError(r.URL.Path)
				if requestID != "" {
					errResp.Meta = map[string]interface{}{
						"request_id": requestID,
					}
				}

				// Write error response
				errors.WriteError(w, errResp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ErrorLogger middleware logs errors for monitoring and debugging
func ErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture status code
		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		// Log errors (4xx and 5xx status codes)
		if wrapper.statusCode >= 400 {
			requestID := GetRequestID(r.Context())
			log.Printf(
				"ERROR: [%s] %s %s - Status: %d - Request-ID: %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				wrapper.statusCode,
				requestID,
			)
		}
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write captures writes (sets status code to 200 if not set)
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}
