package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/direnbharwani/evote-capstone/app/server/common"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
		return errorResponse, nil
	}

	err := common.ChaincodeCastVote(requestBody.VoterID, os.Getenv("KALEIDO_AUTH_TOKEN"), requestBody.BallotID, requestBody.CandidateID)
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("Unable to cast vote: %v", err))
		return errorResponse, nil
	}

	return common.GenerateSuccessResponse(""), nil
}

func main() {
	lambda.Start(Handler)
}

// =============================================================================
// API Types
// =============================================================================

type LambdaRequestBody struct {
	VoterID     string `json:"VoterID"`
	BallotID    string `json:"BallotID"`
	CandidateID string `json:"CandidateID"`
}
