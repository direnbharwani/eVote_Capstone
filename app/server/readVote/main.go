package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/direnbharwani/evote-capstone/app/server/common"
	chaincode "github.com/direnbharwani/evote-capstone/chaincode/src"
	paillier "github.com/direnbharwani/go-paillier/pkg"
)

// ======================================================================================
// Lambda Definition
// ======================================================================================

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("failed to parse request body: %v", err)
	}

	// Invoke Chaincode
	chaincodeResponseBody, err := common.ChaincodeQuery(requestBody.VoterID, "QueryBallot", []string{requestBody.BallotID})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("%v", err)
	}

	ballot := chaincodeResponseBody.Result.(chaincode.Ballot)

	if len(ballot.Candidates) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error: ballot has no candidates")
	}

	// Use private key in conjuction with Candidate's public key
	// to check if candidate has been voted for on a Ballot
	publicKey, privateKey, err := common.DecodeKeys(ballot.Candidates[0].PublicKey, os.Getenv("PAILLIER_PRIVATE_KEY"))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("%v", err)
	}

	responseBody := LambdaResponseBody{BallotID: ballot.Asset.ID}

	for i := range ballot.Candidates {
		decryptedCandidate := LambdaResponseCandidate{
			CandidateID: ballot.Candidates[i].Asset.ID,
			Name:        ballot.Candidates[i].Name,
		}

		encryptedCount, ok := new(big.Int).SetString(ballot.Candidates[i].Count, 10)
		if !ok {
			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("failed to parse candiate count")
		}

		count, err := paillier.Decrypt(publicKey, privateKey, encryptedCount)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error decrypting candidate count: %v", err)
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
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error unparse response body: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(lambdaResponseBodyData),
	}, nil
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
