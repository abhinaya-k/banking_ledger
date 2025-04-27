package models

type ApiError struct {
	StatusCode       int              `json:"statusCode"`
	ApplicationError ApplicationError `json:"applicationError"`
}

type ApplicationError struct {
	Type    string                  `json:"type"`
	Message ApplicationErrorMessage `json:"message"`
}

type ApplicationErrorMessage struct {
	ErrorCode          int         `json:"errorCode"`
	ErrorMessage       string      `json:"errorMessage"`
	DisplayMessage     string      `json:"displayMessage,omitempty"`
	OriginStatusCode   int         `json:"originStatusCode,omitempty"`
	OriginErrorMessage string      `json:"originErrorMessage,omitempty"`
	AdditionalInfo     interface{} `json:"additionalInfo,omitempty"`
}
