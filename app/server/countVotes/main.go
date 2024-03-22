package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	paillier "github.com/direnbharwani/go-paillier/pkg"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("failed to parse request body: %v", err)
	}

	// Invoke chaincode through REST API Gateway to Query Ballot State
	client := &http.Client{}

	chaincodeInvocationHeaders := ChaincodeInvocationHeaders{
		Signer:    requestBody.SignerID,
		Channel:   "default-channel",
		Chaincode: "evote_poc",
	}

	chaincodeRequestBody := ChaincodeRequestBody{
		Headers: chaincodeInvocationHeaders,
		Func:    "QueryAllBallots",
		Args:    []string{},
		Init:    false,
	}

	chaincodeRequestJSONData, err := json.Marshal(chaincodeRequestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("failed to unparse chaincode request body: %v", err)
	}

	chaincodeRequest, err := http.NewRequest("POST", "https://a0z8wc2w78-a0ve7t5vxf-connect.au0-aws-ws.kaleido.io/query", bytes.NewBuffer(chaincodeRequestJSONData))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error creating chaincode request: %v", err)
	}

	chaincodeRequest.Header.Set("Content-Type", "application/json")
	chaincodeRequest.Header.Set("Authorization", os.Getenv("KALEIDO_AUTH_TOKEN"))

	chaincodeResponse, err := client.Do(chaincodeRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error sending chaincode request: %v", err)
	}
	defer chaincodeResponse.Body.Close()

	chaincodeResponseBodyData, err := io.ReadAll(chaincodeResponse.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error chaincode reading response body: %v", err)
	}

	var chaincodeResponseBody ChaincodeQueryResponseBody
	err = json.Unmarshal(chaincodeResponseBodyData, &chaincodeResponseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error parsing chaincode response body: %v", err)
	}

	if len(chaincodeResponseBody.Result) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("no ballots to count")
	}

	// Go through all ballots and store the ones that are part of the election to count votes for
	ballotsToCount := []*ChaincodeBallot{}

	for i := range chaincodeResponseBody.Result {
		if chaincodeResponseBody.Result[i].ElectionID != requestBody.ElectionID {
			continue
		}

		ballotsToCount = append(ballotsToCount, &chaincodeResponseBody.Result[i])
	}

	publicKeyData, err := base64.StdEncoding.DecodeString(chaincodeResponseBody.Result[0].Candidates[0].PublicKey)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error decoding public key: %v", err)
	}

	privateKeyData, err := base64.StdEncoding.DecodeString(os.Getenv("PAILLIER_PRIVATE_KEY"))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error decoding private key: %v", err)
	}

	publicKey, err := paillier.DeserialiseJSON[paillier.PublicKey](publicKeyData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error unparsing public key: %v", err)
	}

	privateKey, err := paillier.DeserialiseJSON[paillier.PrivateKey](privateKeyData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error unparsing private key: %v", err)
	}

	// Create all candidates to count votes for
	results := map[string]*LambdaResponseCandidate{}
	for i := range chaincodeResponseBody.Result[0].Candidates {
		zeroCount, err := paillier.Encrypt(publicKey, big.NewInt(0))
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error encrypting zero count: %v", err)
		}

		candidate := LambdaResponseCandidate{
			CandidateID: chaincodeResponseBody.Result[0].Candidates[i].Asset.ID,
			Name:        chaincodeResponseBody.Result[0].Candidates[i].Name,
			NumVotes:    zeroCount,
		}

		results[candidate.CandidateID] = &candidate
	}

	for i := range ballotsToCount {
		for j := range ballotsToCount[i].Candidates {
			candidate := ballotsToCount[i].Candidates[j]

			candidateCount, ok := new(big.Int).SetString(candidate.Count, 10)
			if !ok {
				return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error parsing candidate count: %v", err)
			}
			currentCount := results[candidate.Asset.ID].NumVotes

			results[candidate.Asset.ID].NumVotes = paillier.AddEncrypted(publicKey, currentCount, candidateCount)
		}
	}

	// Decrypt final count and prepare response
	responseBody := LambdaResponseBody{}

	for _, candidate := range results {
		candidate.NumVotes, err = paillier.Decrypt(publicKey, privateKey, candidate.NumVotes)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error decrypting final count: %v", err)
		}

		responseBody.Candidates = append(responseBody.Candidates, *candidate)
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

// ======================================================================================
// HTTP Types
// ======================================================================================

type LambdaRequestBody struct {
	SignerID   string `json:"SignerID"`
	ElectionID string `json:"ElectionID"`
}

type LambdaResponseBody struct {
	Candidates []LambdaResponseCandidate `json:"Candidates"`
}

type LambdaResponseCandidate struct {
	CandidateID string   `json:"CandidateID"`
	Name        string   `json:"Name"`
	NumVotes    *big.Int `json:"NumVotes"`
}

type ChaincodeInvocationHeaders struct {
	Signer    string `json:"signer"`
	Channel   string `json:"channel"`
	Chaincode string `json:"chaincode"`
}

type ChaincodeRequestBody struct {
	Headers ChaincodeInvocationHeaders `json:"headers"`
	Func    string                     `json:"func"`
	Args    []string                   `json:"args"`
	Init    bool                       `json:"init"`
}

type ChaincodeQueryResponseBody struct {
	Headers map[string]interface{} `json:"headers"`
	Result  []ChaincodeBallot      `json:"result"`
}

// ======================================================================================
// Chaincode Types For Deserialisation
// ======================================================================================

// TODO: There must be a better way to handle this

type Asset struct {
	ID string `json:"ID"`
}

type ChaincodeCandidate struct {
	Asset      Asset  `json:"Asset"`
	Count      string `json:"Count"`
	ElectionID string `json:"ElectionID"`
	Name       string `json:"Name"`
	PublicKey  string `json:"PublicKey"`
}

type ChaincodeBallot struct {
	Asset      Asset                `json:"Asset"`
	Candidates []ChaincodeCandidate `json:"Candidates"`
	ElectionID string               `json:"ElectionID"`
	VoterID    string               `json:"VoterID"`
	Voted      bool                 `json:"Voted"`
}
