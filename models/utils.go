package models

type SuccessResponse struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}
