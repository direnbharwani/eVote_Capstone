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
		Signer:    requestBody.VoterID,
		Channel:   "default-channel",
		Chaincode: "evote_poc",
	}

	chaincodeRequestBody := ChaincodeRequestBody{
		Headers: chaincodeInvocationHeaders,
		Func:    "QueryBallot",
		Args:    []string{requestBody.BallotID},
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

	ballot := chaincodeResponseBody.Result
	responseBody := LambdaResponseBody{BallotID: ballot.Asset.ID}

	if len(ballot.Candidates) != 0 {
		// Use private key in conjuction with Candidate's public key
		// to check if candidate has been voted for on a Ballot
		publicKeyData, err := base64.StdEncoding.DecodeString(ballot.Candidates[0].PublicKey)
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

		for i := range ballot.Candidates {
			var voted bool

			encryptedCount, ok := new(big.Int).SetString(ballot.Candidates[i].Count, 10)
			if !ok {
				return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error setting candidate count: %v", err)
			}

			count, err := paillier.Decrypt(publicKey, privateKey, encryptedCount)
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error decrypting candidate count: %v", err)
			}

			if count.Cmp(big.NewInt(0)) == 0 {
				voted = false
			} else {
				voted = true
			}

			candidate := LambdaResponseCandidate{
				CandidateID: ballot.Candidates[i].Asset.ID,
				Name:        ballot.Candidates[i].Name,
				Voted:       voted,
			}

			responseBody.Candidates = append(responseBody.Candidates, candidate)
		}
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
	Result  ChaincodeBallot        `json:"result"`
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
