package common

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func GenerateSuccessResponse(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"ContentType":                  "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
		},
		Body: body,
	}
}

func GenerateErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	type ErrorBody struct {
		Message string `json:"message"`
	}

	errorResponse := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"ContentType":                  "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
		},
	}

	errorBody, err := json.Marshal(ErrorBody{message})
	if err != nil {
		errorResponse.Body = `{"message": "Bad Request: unable to stringify message!"}`
	} else {
		errorResponse.Body = string(errorBody)
	}

	return errorResponse
}
