package http

import (
	"encoding/json"
	"net/http"
)

type HTTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrorMessageNotFound       = "Resource not found"
	ErrorMessageInternalServer = "Internal Server Error"
	ErrorMessageBadRequest     = "Bad Request"
	ErrorMessageUnauthorized   = "Unauthorized"
	ErrorMessageForbidden      = "Forbidden"
)

type ResponseCfg struct {
	Code    int
	Data    any
	Headers map[string]string
	Err     error
}

func WriteJsonResponse(w http.ResponseWriter, config *ResponseCfg) {
	w.Header().Set("Content-Type", "application/json")
	for key, value := range config.Headers {
		w.Header().Add(key, value)
	}

	if config.Err != nil {
		WriteError(w, config)
		return
	}
	w.WriteHeader(config.Code)
	if err := json.NewEncoder(w).Encode(config.Data); err != nil {
		http.Error(
			w,
			ErrorMessageInternalServer,
			http.StatusInternalServerError,
		)
		return
	}
}

func WriteError(w http.ResponseWriter, config *ResponseCfg) {
	w.Header().Set("Content-Type", "application/json")
	for key, value := range config.Headers {
		w.Header().Add(key, value)
	}

	errorData := HTTTPError{
		Code:    config.Code,
		Message: getErrorMessage(config.Code),
	}

	w.WriteHeader(config.Code)
	if err := json.NewEncoder(w).Encode(errorData); err != nil {
		http.Error(
			w,
			ErrorMessageInternalServer,
			http.StatusInternalServerError,
		)
		return
	}
}

func getErrorMessage(code int) string {
	switch code {
	case http.StatusNotFound:
		return ErrorMessageNotFound
	case http.StatusBadRequest:
		return ErrorMessageBadRequest
	case http.StatusUnauthorized:
		return ErrorMessageUnauthorized
	case http.StatusForbidden:
		return ErrorMessageForbidden
	default:
		return ErrorMessageInternalServer
	}
}

func SuccessResponse(
	w http.ResponseWriter,
	code int,
	data any,
	headers map[string]string,
) {
	WriteJsonResponse(w, &ResponseCfg{
		Code:    code,
		Data:    data,
		Headers: headers,
	})
}

func ErrorResponse(
	w http.ResponseWriter,
	code int,
	err error,
	headers map[string]string,
) {
	WriteError(w, &ResponseCfg{
		Code:    code,
		Err:     err,
		Headers: headers,
	})
}
