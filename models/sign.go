package models

type SignRequest struct {
	Keys []string `json:"keys" validate:"required,min=1"`
}

type SignResponse struct {
	URLs map[string]string `json:"urls"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
