package apiresponse

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func write(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(response)
}

func Success(
	w http.ResponseWriter,
	status int,
	message string,
	data any,
) {
	write(w, status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(
	w http.ResponseWriter,
	status int,
	message string,
	data any,
	pagination Pagination,
) {
	write(w, status, APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: &pagination,
	})
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
) {
	write(w, status, APIResponse{
		Success: false,
		Message: message,
	})
}

func OK(w http.ResponseWriter, data any) {
	Success(
		w,
		http.StatusOK,
		"Success",
		data,
	)
}

func Created(w http.ResponseWriter, data any) {
	Success(
		w,
		http.StatusCreated,
		"Created successfully",
		data,
	)
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message)
}

func InternalServerError(w http.ResponseWriter) {
	Error(
		w,
		http.StatusInternalServerError,
		"Internal server error",
	)
}
