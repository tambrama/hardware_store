package dto

type ErrorResponse struct {
	Error string `json:"error" example:"validation error"`
}

type ValidationErrorResponse struct {
	Error string `json:"error" example:"invalid uuid format"`
}

type NotFoundErrorResponse struct {
	Error string `json:"error" example:"client not found"`
}

type InternalErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
