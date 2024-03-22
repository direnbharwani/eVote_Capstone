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
	paillier "github.com/direnbharwani/evote-capstone/paillier"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("failed to parse request body: %v", err)
	}

	ballots, err := common.ChaincodeQueryAll[chaincode.Ballot](requestBody.SignerID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("%v", err)
	}

	if len(ballots) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("no ballots to count")
	}

	// Filter ballots that are part of specified election
	ballotsToCount := []chaincode.Ballot{}

	for i := range ballots {
		if ballots[i].ElectionID != requestBody.ElectionID {
			ballotsToCount = append(ballotsToCount, ballots[i])
		}
	}
	if len(ballotsToCount) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("no ballots to count")
	}
	if len(ballotsToCount[0].Candidates) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("ballot does not have candidates")
	}

	// First publicKey is used for comparison against other public keys as a control measure
	// Privat key will be used at the end for decypting the final count
	publicKey, privateKey, err := common.DecodeKeys(ballotsToCount[0].Candidates[0].PublicKey, os.Getenv("PAILLIER_PRIVATE_KEY"))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("%v", err)
	}

	// Create all candidates to count votes for
	results, err := countBallots(ballotsToCount, publicKey)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("%v", err)
	}

	// Decrypt final count and prepare response
	for i := range results {
		results[i].NumVotes, err = paillier.Decrypt(publicKey, privateKey, results[i].NumVotes)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error decrypting final count: %v", err)
		}
	}

	lambdaResponseBodyData, err := json.Marshal(LambdaResponseBody{results})
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
// Helper Methods
// ======================================================================================

func countBallots(ballotsToCount []chaincode.Ballot, publicKey *paillier.PublicKey) ([]LambdaResponseCandidate, error) {
	// Create set of candidates to count in a map
	// We use a map since access is faster to update: O(1)
	candidateMap := map[string]LambdaResponseCandidate{}

	// We take the first ballot's candidates
	// This will be used as a control measure to ensure all ballots have at least one of these candidates
	// The onus of ensuring ballots have all the candidates (and no more) is not within the scope of this lambda
	zeroCount, err := paillier.Encrypt(publicKey, big.NewInt(0))
	if err != nil {
		return []LambdaResponseCandidate{}, err
	}

	for i := range ballotsToCount[0].Candidates {
		candidate := ballotsToCount[0].Candidates[i]

		candidateMap[candidate.Asset.ID] = LambdaResponseCandidate{
			CandidateID: candidate.Asset.ID,
			Name:        candidate.Name,
			NumVotes:    zeroCount,
		}
	}

	// Go through all candidates in each ballot and add the count to the candidate in the map
	// Unavoidable nested loop: O(nc) where n is number of ballots, c is number of candidates
	for i := range ballotsToCount {
		ballot := ballotsToCount[i]

		for j := range ballot.Candidates {
			candidate := ballot.Candidates[j]

			if c, found := candidateMap[candidate.Asset.ID]; found {
				candidateCount, ok := new(big.Int).SetString(candidate.Count, 10)
				if !ok {
					return []LambdaResponseCandidate{}, fmt.Errorf("error parsing candidate count: %v", err)
				}

				// Replace candiadte in map
				candidateMap[candidate.Asset.ID] = LambdaResponseCandidate{
					CandidateID: candidate.Asset.ID,
					Name:        candidate.Name,
					NumVotes:    paillier.AddEncrypted(publicKey, c.NumVotes, candidateCount),
				}
			} else {
				return []LambdaResponseCandidate{}, fmt.Errorf("extra candidate found in ballot %s", ballot.Asset.ID)
			}
		}
	}

	// Create final array with candidates from the map
	results := []LambdaResponseCandidate{}
	for _, c := range candidateMap {
		results = append(results, c)
	}

	return results, nil
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
