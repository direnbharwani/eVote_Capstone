package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Function to test if the chaincode has been successfully deployed
func (s *SmartContract) LiveTest() string {
	data := map[string]interface{}{
		"Name":    "Capstone eVote POC Chaincode",
		"Version": "v0.1",
		"Time":    time.Now().Format(time.DateTime),
		"Status":  "Live",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	}

	return string(jsonData)
}

// =============================================================================
// Creation
// =============================================================================

// Creates a ballot as an asset on the blockchain
// data must contain BallotID, VoterID & ElectionID
func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, data string) error {
	ballot, err := ParseJSON[Ballot](data)
	if err != nil {
		return err
	}

	election, err := queryAsset[Election](ctx, ballot.ElectionID)
	if err != nil {
		return err
	}

	ballot.Candidates = election.Candidates
	ballot.Voted = false

	return createAsset(ctx, ballot.BallotID, ballot, "ballot")
}

// Creates a candidate as an asset on the blockchain
// data must contain CandidateID & ElectionID
func (s *SmartContract) CreateCandidate(ctx contractapi.TransactionContextInterface, data string) error {
	candidate, err := ParseJSON[Candidate](data)
	if err != nil {
		return err
	}

	candidate.Count = 0
	if err := createAsset(ctx, candidate.CandidateID, candidate, "candidate"); err != nil {
		return err
	}

	// Update election with candidate
	election, err := queryAsset[Election](ctx, candidate.ElectionID)
	if err != nil {
		return err
	}

	election.Candidates = append(election.Candidates, candidate)
	return updateAsset(ctx, election.ElectionID, election)
}

func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, data string) error {
	election, err := ParseJSON[Election](data)
	if err != nil {
		return err
	}

	// TODO: Sync by fetching all candidates and taking ones with matching electionIDs

	return createAsset(ctx, election.ElectionID, election, "election")
}

func (s *SmartContract) CreateVoter(ctx contractapi.TransactionContextInterface, data string) error {
	voter, err := ParseJSON[Voter](data)
	if err != nil {
		return err
	}

	return createAsset(ctx, voter.VoterID, voter, "voter")
}

func createAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string, createdAsset T, assetType string) error {
	objectState, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}
	if objectState != nil {
		return fmt.Errorf("%s: %s already created", assetType, key)
	}

	createdData, err := json.Marshal(createdAsset)
	if err != nil {
		return err
	}

	if err = ctx.GetStub().PutState(key, createdData); err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}

	return nil
}

// =============================================================================
// Query
// =============================================================================

func (s *SmartContract) QueryBallot(ctx contractapi.TransactionContextInterface, key string) (Ballot, error) {
	return queryAsset[Ballot](ctx, key)
}

func (s *SmartContract) QueryCandidate(ctx contractapi.TransactionContextInterface, key string) (Candidate, error) {
	return queryAsset[Candidate](ctx, key)
}

func (s *SmartContract) QueryElection(ctx contractapi.TransactionContextInterface, key string) (Election, error) {
	return queryAsset[Election](ctx, key)
}

func (s *SmartContract) QueryVoter(ctx contractapi.TransactionContextInterface, key string) (Voter, error) {
	return queryAsset[Voter](ctx, key)
}

func (s *SmartContract) QueryAllBallots(ctx contractapi.TransactionContextInterface) ([]Ballot, error) {
	return queryAssetsByType[Ballot](ctx)
}

func (s *SmartContract) QueryAllCandidates(ctx contractapi.TransactionContextInterface) ([]Candidate, error) {
	return queryAssetsByType[Candidate](ctx)
}

func (s *SmartContract) QueryAllElections(ctx contractapi.TransactionContextInterface) ([]Election, error) {
	return queryAssetsByType[Election](ctx)
}

func (s *SmartContract) QueryAllVoters(ctx contractapi.TransactionContextInterface) ([]Voter, error) {
	return queryAssetsByType[Voter](ctx)
}

func queryAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string) (T, error) {
	var emptyObject T
	var result T

	assetState, err := ctx.GetStub().GetState(key)
	if err != nil {
		return emptyObject, fmt.Errorf("unable to interact with world state: %v", err)
	}
	if assetState == nil {
		return emptyObject, fmt.Errorf("cannot read world state with key %s", key)
	}

	if err = json.Unmarshal(assetState, &result); err != nil {
		return emptyObject, err
	}

	return result, nil
}

func queryAssetsByType[T ITYPES](ctx contractapi.TransactionContextInterface) ([]T, error) {
	// By using empty start & end keys, we will grab all assets
	startKey := ""
	endKey := ""

	results := []T{}
	resultIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}

	defer resultIterator.Close()
	for resultIterator.HasNext() {
		var result T
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(queryResponse.Value, &result); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

// =============================================================================
// Update
// =============================================================================

// Updates a ballot with the specified updated state.
// The ballot cannot be updated if the ballot has already been cast.
func (s *SmartContract) UpdateBallot(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedBallot, err := ParseJSON[Ballot](updatedData)
	if err != nil {
		return err
	}

	currentBallot, err := queryAsset[Ballot](ctx, updatedBallot.BallotID)
	if err != nil {
		return err
	}

	if currentBallot.Voted {
		return fmt.Errorf("unable to update ballot %s that has already been voted", currentBallot.BallotID)
	}

	return updateAsset(ctx, updatedBallot.BallotID, updatedBallot)
}

func (s *SmartContract) UpdateCandidate(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedCandidate, err := ParseJSON[Candidate](updatedData)
	if err != nil {
		return err
	}

	return updateAsset(ctx, updatedCandidate.CandidateID, updatedCandidate)
}

func (s *SmartContract) UpdateElection(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedElection, err := ParseJSON[Election](updatedData)
	if err != nil {
		return err
	}

	return updateAsset(ctx, updatedElection.ElectionID, updatedElection)
}

func (s *SmartContract) UpdateVoter(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedVoter, err := ParseJSON[Voter](updatedData)
	if err != nil {
		return err
	}

	return updateAsset(ctx, updatedVoter.VoterID, updatedVoter)
}

func updateAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string, updatedAsset T) error {
	assetState, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}
	if assetState == nil {
		return fmt.Errorf("cannot read world state with key %s", key)
	}

	bytes, err := json.Marshal(updatedAsset)
	if err != nil {
		return err
	}

	if err = ctx.GetStub().PutState(key, bytes); err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}

	return nil
}

// =============================================================================
// Update
// =============================================================================

func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterID string, ballotID string, candidateID string) error {
	voter, err := queryAsset[Voter](ctx, voterID)
	if err != nil {
		return err
	}

	ballot, err := queryAsset[Ballot](ctx, ballotID)
	if err != nil {
		return err
	}

	if voter.BallotID != ballotID || ballot.VoterID != voterID {
		errorMessage := fmt.Sprintf("voter %s is not assigned ballot %s!", voterID, ballotID)
		return errors.New(errorMessage)
	}

	candidateFound := false
	for i, c := range ballot.Candidates {
		if c.CandidateID == candidateID {
			ballot.Candidates[i].Count++
			candidateFound = true
			break
		}
	}

	if !candidateFound {
		errorMessage := fmt.Sprintf("candidate %s is not found in ballot %s!", candidateID, ballotID)
		return errors.New(errorMessage)
	}

	return updateAsset[Ballot](ctx, ballot.BallotID, ballot)
}