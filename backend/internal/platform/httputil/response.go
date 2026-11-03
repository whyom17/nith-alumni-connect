package httputil

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the standard API error envelope.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail holds details about a specific error.
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// WriteJSON serializes payload as JSON and writes it to the response writer.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// WriteError writes a standardized error envelope JSON response.
func WriteError(w http.ResponseWriter, status int, code, message string, fields map[string]string) {
	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	}
	WriteJSON(w, status, resp)
}
