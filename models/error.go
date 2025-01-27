package models

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewErrorResponse(error string, message string) ErrorResponse {
	return ErrorResponse{
		Error:   error,
		Message: message,
	}
}
