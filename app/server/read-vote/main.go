package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/direnbharwani/evote-capstone/app/server/common"
	chaincode "github.com/direnbharwani/evote-capstone/chaincode/src"
	paillier "github.com/direnbharwani/evote-capstone/paillier"
)

// ======================================================================================
// Lambda Definition
// ======================================================================================

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
		return errorResponse, nil
	}

	// Invoke Chaincode
	ballot, err := common.ChaincodeQuery[chaincode.Ballot](requestBody.VoterID, os.Getenv("KALEIDO_AUTH_TOKEN"), requestBody.BallotID)
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	if len(ballot.Candidates) == 0 {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, "error: ballot has no candidates")
		return errorResponse, nil
	}

	// Use private key in conjuction with Candidate's public key
	// to check if candidate has been voted for on a Ballot
	publicKey, privateKey, err := common.DecodeKeys(ballot.Candidates[0].PublicKey, os.Getenv("PAILLIER_PRIVATE_KEY"))
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	responseBody := LambdaResponseBody{BallotID: ballot.Asset.ID}

	for i := range ballot.Candidates {
		decryptedCandidate := LambdaResponseCandidate{
			CandidateID: ballot.Candidates[i].Asset.ID,
			Name:        ballot.Candidates[i].Name,
		}

		encryptedCount, ok := new(big.Int).SetString(ballot.Candidates[i].Count, 10)
		if !ok {
			errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
			return errorResponse, nil
		}

		count, err := paillier.Decrypt(publicKey, privateKey, encryptedCount)
		if err != nil {
			errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
			return errorResponse, nil
		}

		if count.Cmp(big.NewInt(0)) == 0 {
			decryptedCandidate.Voted = false
		} else {
			decryptedCandidate.Voted = true
		}

		responseBody.Candidates = append(responseBody.Candidates, decryptedCandidate)
	}

	lambdaResponseBodyData, err := json.Marshal(responseBody)
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("error unparse response body: %v", err))
		return errorResponse, nil
	}

	return common.GenerateSuccessResponse(string(lambdaResponseBodyData)), nil
}

func main() {
	lambda.Start(Handler)
}

// =============================================================================
// API Types
// =============================================================================

type LambdaRequestBody struct {
	VoterID    string `json:"VoterID"`
	BallotID   string `json:"BallotID"`
	ElectionID string `json:"ElectionID"`
}

type LambdaResponseBody struct {
	BallotID   string                    `json:"BallotID"`
	Candidates []LambdaResponseCandidate `json:"Candidates"`
}

type LambdaResponseCandidate struct {
	CandidateID string `json:"CandidateID"`
	Name        string `json:"Name"`
	Voted       bool   `json:"Voted"`
}
