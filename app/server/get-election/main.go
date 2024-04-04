package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/direnbharwani/evote-capstone/app/server/common"
	chaincode "github.com/direnbharwani/evote-capstone/chaincode/src"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract the qualifier from the path parameters
	electionID := request.PathParameters["qualifier"]
	fmt.Printf("Received qualifier: %s\n", electionID)

	election, err := common.ChaincodeQuery[chaincode.Election]("testVoter0", os.Getenv("KALEIDO_AUTH_TOKEN"), electionID)
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	lambdaResponseBody := LambdaResponseBody{
		Election: election,
		IsActive: election.IsActive(),
	}

	lambdaResponseBodyData, err := json.Marshal(lambdaResponseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return common.GenerateSuccessResponse(string(lambdaResponseBodyData)), nil
}

func main() {
	lambda.Start(handler)
}

// =============================================================================
// API Types
// =============================================================================

type LambdaResponseBody struct {
	Election chaincode.Election `json:"Election"`
	IsActive bool               `json:"isActive`
}
