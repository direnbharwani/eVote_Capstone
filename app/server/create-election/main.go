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
	"github.com/google/uuid"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
		return errorResponse, nil
	}

	// Create Election

	electionID, err := uuid.NewV7()
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("error generating election id: %v", err))
		return errorResponse, nil
	}

	newElection := chaincode.Election{
		Asset:     chaincode.Asset{ID: "e-" + electionID.String()},
		EndTime:   requestBody.EndTime,
		Name:      requestBody.ElectionName,
		StartTime: requestBody.StartTime,
	}

	if err = common.ChaincodeCreate("testVoter0", os.Getenv("KALEIDO_AUTH_TOKEN"), newElection); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	// Create n number of candidates
	for i := 0; i < requestBody.NumCandidates; i++ {
		candidateID, err := uuid.NewV7()
		if err != nil {
			errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("error generating election id: %v", err))
			return errorResponse, nil
		}

		newCandidate := chaincode.Candidate{
			Asset:      chaincode.Asset{ID: "c-" + candidateID.String()},
			ElectionID: newElection.Asset.ID,
			Name:       fmt.Sprintf("candidate-%d", i),
			PublicKey:  os.Getenv("PAILLIER_PUBLIC_KEY"),
		}

		if err = common.ChaincodeCreate("testVoter0", os.Getenv("KALEIDO_AUTH_TOKEN"), newCandidate); err != nil {
			errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
			return errorResponse, nil
		}
	}

	// Sync
	if err = common.ChaincodeSync("testVoter0", os.Getenv("KALEIDO_AUTH_TOKEN"), newElection.Asset.ID); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	// Return electionID
	responseBody := LambdaResponseBody{newElection.Asset.ID}
	lambdaResponseBodyData, err := json.Marshal(responseBody)
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("error unparse response body: %v", err))
		return errorResponse, nil
	}

	return common.GenerateSuccessResponse(string(lambdaResponseBodyData)), nil
}

func main() {
	lambda.Start(handler)
}

// ======================================================================================
// HTTP Types
// ======================================================================================

type LambdaRequestBody struct {
	ElectionName  string `json:"ElectionName"`
	StartTime     string `json:"StartTime"`
	EndTime       string `json:"EndTime"`
	NumCandidates int    `json:"NumCandidates"`
}

type LambdaResponseBody struct {
	ElectionID string `json:"ElectionID"`
}
