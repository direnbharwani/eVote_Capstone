package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ITYPES is a union set type constraint
// that enforces only allowable types are passed to smart contract methods.
type ITYPES interface {
	Ballot | Candidate | Election | Voter
}

type SmartContract struct {
	contractapi.Contract
}

// Function to test if the chaincode has been successfully deployed
func (s *SmartContract) LiveTest() string {
	return "Hello EVote V1!"
}

// =============================================================================
// Creation
// =============================================================================

func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, voterID string, electionID string) error {
	// Check if voter has an existing ballot
	voter, err := queryAsset[Voter](ctx, voterID)
	if err != nil {
		return err
	}

	if voter.BallotID != "" { // existing ballot
		return fmt.Errorf("voter %s has an existing ballot %s", voterID, voter.BallotID)
	}

	// Get election state
	election, err := queryAsset[Election](ctx, electionID)
	if err != nil {
		return err
	}

	// Create Ballot
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	newBallot := Ballot{
		BallotID:   "b-" + uuid.String(),
		Candidates: election.Candidates,
		ElectionID: electionID,
		VoterID:    voterID,
		Voted:      false,
	}

	// Update Voter with new BallotID
	voter.BallotID = newBallot.BallotID

	err = createAsset(ctx, newBallot.BallotID, newBallot, "ballot")
	if err != nil {
		return err
	}

	err = updateAsset(ctx, voterID, voter)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) CreateCandidate(ctx contractapi.TransactionContextInterface, candidate Candidate) error {
	return createAsset[Candidate](ctx, candidate.CandidateID, candidate, "candidate")
}

func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, election Election) error {
	return createAsset[Election](ctx, election.ElectionID, election, "election")
}

func (s *SmartContract) CreateVoter(ctx contractapi.TransactionContextInterface, voter Voter) error {
	return createAsset[Voter](ctx, voter.VoterID, voter, "voter")
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

	err = ctx.GetStub().PutState(key, createdData)
	if err != nil {
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

	err = json.Unmarshal(assetState, &result)
	if err != nil {
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

		err = json.Unmarshal(queryResponse.Value, &result)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

// =============================================================================
// Update
// =============================================================================

func (s *SmartContract) UpdateBallot(ctx contractapi.TransactionContextInterface, updatedBallot Ballot) error {
	currentBallot, err := queryAsset[Ballot](ctx, updatedBallot.BallotID)
	if err != nil {
		return err
	}

	if currentBallot.Voted {
		return fmt.Errorf("unable to update ballot %s that has already been voted", currentBallot.BallotID)
	}

	updatedBallot.Voted = true
	return updateAsset[Ballot](ctx, updatedBallot.BallotID, updatedBallot)
}

func (s *SmartContract) UpdateCandidate(ctx contractapi.TransactionContextInterface, updatedCandidate Candidate) error {
	return updateAsset[Candidate](ctx, updatedCandidate.CandidateID, updatedCandidate)
}

func (s *SmartContract) UpdateElection(ctx contractapi.TransactionContextInterface, updatedElection Election) error {
	return updateAsset[Election](ctx, updatedElection.ElectionID, updatedElection)
}

func (s *SmartContract) UpdateVoter(ctx contractapi.TransactionContextInterface, updatedVoter Voter) error {
	return updateAsset[Voter](ctx, updatedVoter.VoterID, updatedVoter)
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

	err = ctx.GetStub().PutState(key, bytes)
	if err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}

	return nil
}
