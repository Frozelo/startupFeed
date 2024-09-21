package http

import (
	"encoding/json"
	"net/http"
)

// Стандартные сообщения об ошибках
const (
	ErrorMessageNotFound       = "Resource not found"
	ErrorMessageInternalServer = "Internal Server Error"
	ErrorMessageBadRequest     = "Bad Request"
	ErrorMessageUnauthorized   = "Unauthorized"
	ErrorMessageForbidden      = "Forbidden"
)

// Структура для конфигурации ошибок с подробной информацией
type ErrorResponse struct {
	Code    int    `json:"code"`              // Код ошибки
	Message string `json:"message"`           // Основное сообщение об ошибке
	Details string `json:"details,omitempty"` // Дополнительные подробности об ошибке
}

// Структура для успешного ответа
type SuccessResponse struct {
	Code int `json:"code"`
	Data any `json:"data,omitempty"`
}

// Структура для отправки ответа
type Response struct {
	Code    int               `json:"code"`
	Data    any               `json:"data,omitempty"`
	Message string            `json:"message,omitempty"`
	Details string            `json:"details,omitempty"`
	Headers map[string]string `json:"-"`
}

// Функция для отправки JSON-ответа
func WriteJSONResponse(w http.ResponseWriter, config *Response) {
	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	for key, value := range config.Headers {
		w.Header().Set(key, value)
	}

	// Устанавливаем код статуса
	w.WriteHeader(config.Code)

	// Пишем JSON-ответ
	if err := json.NewEncoder(w).Encode(config); err != nil {
		// Если возникла ошибка сериализации, отправляем внутреннюю ошибку сервера
		http.Error(
			w,
			ErrorMessageInternalServer,
			http.StatusInternalServerError,
		)
	}
}

// Функция для успешного ответа
func Success(
	w http.ResponseWriter,
	code int,
	data any,
	headers map[string]string,
) {
	response := &Response{
		Code:    code,
		Data:    data,
		Headers: headers,
	}
	WriteJSONResponse(w, response)
}

// Функция для отправки ошибки
func Error(
	w http.ResponseWriter,
	code int,
	err error,
	details string,
	headers map[string]string,
) {
	response := &Response{
		Code:    code,
		Message: getErrorMessage(code),
		Details: details,
		Headers: headers,
	}
	if err != nil {
		response.Message = err.Error() // Используем сообщение из самой ошибки
	}
	WriteJSONResponse(w, response)
}

// Получаем стандартное сообщение об ошибке по коду
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
