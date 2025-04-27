package models

const (
	MISC_ERROR                           = 1
	KAFKA_ERROR_NO_INTERVENTION_REQUIRED = 2
	API_ERROR_NO_INTERVENTION_REQUIRED   = 3
	ERROR_NO_INTERVENTION_REQUIRED       = 4
	ERROR_REQUIRE_INTERVENTION           = 5
	KAFKA_ERROR_REQUIRE_INTERVENTION     = 6
	API_ERROR_REQUIRE_INTERVENTION       = 7
	KAFKA_PRODUCER_ERROR                 = 8
	KAFKA_CONSUMER_ERROR                 = 9
)

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

type DroppedMessage struct {
	TopicName    string `json:"topicName"`
	ErrorType    string `json:"errorType"`
	KafkaMessage string `json:"kafkaMessage"`
}
