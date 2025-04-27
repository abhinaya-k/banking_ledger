package models

var CONTEXT_REQUEST_ID_KEY string = "requestId"

type SuccessResponse struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}
