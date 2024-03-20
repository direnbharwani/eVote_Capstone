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
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return err.Error()
	}

	data := map[string]interface{}{
		"Name":    "Capstone eVote POC Chaincode",
		"Version": "v1.2",
		"Time":    time.Now().In(loc).Format(time.DateTime),
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
// data must contain Asset.ID, VoterID & ElectionID
func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, data string) error {
	ballot, err := ParseJSON[Ballot](data)
	if err != nil {
		return err
	}

	election, err := queryAsset[Election](ctx, ballot.ElectionID)
	if err != nil {
		return err
	}

	for _, candidateID := range election.Candidates {
		candidate, err := queryAsset[Candidate](ctx, candidateID)
		if err != nil {
			return err
		}

		ballot.Candidates = append(ballot.Candidates, candidate)
	}

	// Default state must be false
	ballot.Voted = false

	return createAsset(ctx, ballot.Asset.ID, ballot)
}

// Creates a candidate as an asset on the blockchain
// data must contain Asset.ID & ElectionID
func (s *SmartContract) CreateCandidate(ctx contractapi.TransactionContextInterface, data string) error {
	candidate, err := ParseJSON[Candidate](data)
	if err != nil {
		return err
	}

	// Default state must be 0 count. Count will not change on candidate assets, only in ballots.
	candidate.Count = 0
	return createAsset(ctx, candidate.Asset.ID, candidate)
}

func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, data string) error {
	election, err := ParseJSON[Election](data)
	if err != nil {
		return err
	}

	// TODO: Sync by fetching all candidates and taking ones with matching electionIDs

	return createAsset(ctx, election.Asset.ID, election)
}

func createAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string, createdAsset T) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey(createdAsset.Type(), []string{key})
	if err != nil {
		return &CompositeKeyCreationError{err.Error(), key, createdAsset.Type()}
	}

	assetState, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return &WorldStateInteractionError{key, err.Error()}
	}
	if assetState != nil {
		return fmt.Errorf("%s: %s already created", createdAsset.Type(), key)
	}

	createdData, err := json.Marshal(createdAsset)
	if err != nil {
		return err
	}

	if err = ctx.GetStub().PutState(compositeKey, createdData); err != nil {
		return &WorldStateInteractionError{key, err.Error()}
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

func (s *SmartContract) QueryAllBallots(ctx contractapi.TransactionContextInterface) ([]Ballot, error) {
	return queryAssetsByType[Ballot](ctx)
}

func (s *SmartContract) QueryAllCandidates(ctx contractapi.TransactionContextInterface) ([]Candidate, error) {
	return queryAssetsByType[Candidate](ctx)
}

func (s *SmartContract) QueryAllElections(ctx contractapi.TransactionContextInterface) ([]Election, error) {
	return queryAssetsByType[Election](ctx)
}

func queryAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string) (T, error) {
	var emptyObject T
	var result T

	compositeKey, err := ctx.GetStub().CreateCompositeKey(result.Type(), []string{key})
	if err != nil {
		return emptyObject, &CompositeKeyCreationError{err.Error(), key, result.Type()}
	}

	assetState, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return emptyObject, &WorldStateInteractionError{err.Error(), key}
	}
	if assetState == nil {
		return emptyObject, &WorldStateReadFailureError{key}
	}

	if err = json.Unmarshal(assetState, &result); err != nil {
		return emptyObject, err
	}

	return result, nil
}

func queryAssetsByType[T ITYPES](ctx contractapi.TransactionContextInterface) ([]T, error) {
	var emptyObject T

	results := []T{}
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(emptyObject.Type(), []string{})
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	for resultIterator.HasNext() {
		compositeKey, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		_, keys, err := ctx.GetStub().SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return nil, err
		}

		result, err := queryAsset[T](ctx, keys[0])
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

// Updates a ballot with the specified updated state.
// The ballot cannot be updated if the ballot has already been cast.
func (s *SmartContract) UpdateBallot(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedState, err := ParseJSON[Ballot](updatedData)
	if err != nil {
		return err
	}

	currentState, err := queryAsset[Ballot](ctx, updatedState.Asset.ID)
	if err != nil {
		return err
	}

	if currentState.Voted {
		return fmt.Errorf("unable to update ballot %s that has already been voted", currentState.Asset.ID)
	}

	return updateAsset(ctx, updatedState.Asset.ID, updatedState)
}

func (s *SmartContract) UpdateCandidate(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedState, err := ParseJSON[Candidate](updatedData)
	if err != nil {
		return err
	}

	return updateAsset(ctx, updatedState.Asset.ID, updatedState)
}

func (s *SmartContract) UpdateElection(ctx contractapi.TransactionContextInterface, updatedData string) error {
	updatedState, err := ParseJSON[Election](updatedData)
	if err != nil {
		return err
	}

	return updateAsset(ctx, updatedState.Asset.ID, updatedState)
}

func updateAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string, updatedAsset T) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey(updatedAsset.Type(), []string{key})
	if err != nil {
		return &CompositeKeyCreationError{err.Error(), key, updatedAsset.Type()}
	}

	currentState, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return &WorldStateInteractionError{err.Error(), key}
	}
	if currentState == nil {
		return &WorldStateReadFailureError{key}
	}

	var currentAsset T
	if err = json.Unmarshal(currentState, &currentAsset); err != nil {
		return err
	}

	if currentAsset.IsEqual(updatedAsset) {
		return &ObjectEqualityError{key, updatedAsset.Type()}
	}

	updatedData, err := json.Marshal(updatedAsset)
	if err != nil {
		return err
	}

	if err = ctx.GetStub().PutState(compositeKey, updatedData); err != nil {
		return &WorldStateInteractionError{err.Error(), key}
	}

	return nil
}

// =============================================================================
// Deletion (only for testing)
// =============================================================================

func (s *SmartContract) DeleteElection(ctx contractapi.TransactionContextInterface, key string) error {
	return deleteAsset[Election](ctx, key)
}

func (s *SmartContract) DeleteCandidate(ctx contractapi.TransactionContextInterface, key string) error {
	return deleteAsset[Candidate](ctx, key)
}

func (s *SmartContract) DeleteBallot(ctx contractapi.TransactionContextInterface, key string) error {
	return deleteAsset[Ballot](ctx, key)
}

func deleteAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey(updatedAsset.Type(), []string{key})
	if err != nil {
		return &CompositeKeyCreationError{err.Error(), key, updatedAsset.Type()}
	}

	currentState, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return &WorldStateInteractionError{err.Error(), key}
	}
	if currentState == nil {
		return &WorldStateReadFailureError{key}
	}

	if err := ctx.GetStub().DelState(compositeKey); err != nil {
		return err
	}

	return nil
}

// =============================================================================
// Cast Vote
// =============================================================================

func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterID string, ballotID string, candidateID string) error {
	ballot, err := queryAsset[Ballot](ctx, ballotID)
	if err != nil {
		return err
	}

	if ballot.VoterID != voterID {
		errorMessage := fmt.Sprintf("voter %s is not assigned ballot %s!", voterID, ballotID)
		return errors.New(errorMessage)
	}

	// Ensure election is active
	election, err := queryAsset[Election](ctx, ballot.ElectionID)
	if err != nil {
		return err
	}
	if !election.IsActive() {
		errorMessage := fmt.Sprintf("election %s is not active! vote cannot be cast", election.Asset.ID)
		return errors.New(errorMessage)
	}

	if err := ballot.Vote(candidateID); err != nil {
		return err
	}

	return updateAsset[Ballot](ctx, ballot.Asset.ID, ballot)
}
