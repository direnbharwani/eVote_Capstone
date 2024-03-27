package chaincode

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
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
		"Name":    "eVote POC Chaincode",
		"Version": "v2.4.1-performance-test",
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
// data must contain Asset.ID & ElectionID
// No candidates are expected as they will be taken from the Election asset
func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, data string) error {
	ballot, err := ParseJSON[Ballot](data)
	if err != nil {
		return err
	}

	election, err := queryAsset[Election](ctx, ballot.ElectionID)
	if err != nil {
		return err
	}

	// Clear the ballot slice to ensure no duplicate candides
	ballot.Candidates = ballot.Candidates[:0]

	for _, candidateID := range election.Candidates {
		candidate, err := queryAsset[Candidate](ctx, candidateID)
		if err != nil {
			return err
		}

		if err = candidate.Init(); err != nil {
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
	candidate.Init()
	return createAsset(ctx, candidate.Asset.ID, candidate)
}

// Creates an election as an asset on the blockchain
// data must contian Asset.ID, StartTime & EndTime.
// StartTime must be before EndTime.
func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, data string) error {
	election, err := ParseJSON[Election](data)
	if err != nil {
		return err
	}

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

func (s *SmartContract) QueryBallotHistory(ctx contractapi.TransactionContextInterface, key string) (map[string]Ballot, error) {
	return queryAssetHistory[Ballot](ctx, key)
}

func (s *SmartContract) QueryCandidateHistory(ctx contractapi.TransactionContextInterface, key string) (map[string]Candidate, error) {
	return queryAssetHistory[Candidate](ctx, key)
}

func (s *SmartContract) QueryElectionHistory(ctx contractapi.TransactionContextInterface, key string) (map[string]Election, error) {
	return queryAssetHistory[Election](ctx, key)
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

func queryAssetHistory[T ITYPES](ctx contractapi.TransactionContextInterface, key string) (map[string]T, error) {
	var emptyObject T

	compositeKey, err := ctx.GetStub().CreateCompositeKey(emptyObject.Type(), []string{key})
	if err != nil {
		return nil, &CompositeKeyCreationError{err.Error(), key, emptyObject.Type()}
	}

	assetHistory := make(map[string]T)
	resultIterator, err := ctx.GetStub().GetHistoryForKey(compositeKey)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	// We continue on errors to avoid returning an error due to deleted assets
	for resultIterator.HasNext() {
		assetState, err := resultIterator.Next()
		if err != nil {
			fmt.Printf("failed to retrieve state for %s %s\n", emptyObject.Type(), key)
			continue
		}

		var result T
		if err = json.Unmarshal(assetState.Value, &result); err != nil {
			fmt.Printf("failed to parse state for %s %s\n", emptyObject.Type(), key)
			continue
		}

		assetHistory[assetState.TxId] = result
	}

	// TODO: Sort the transactions by their timestamp, from earliest (at 0) to latest (at len-1)

	return assetHistory, nil
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
		assetState, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var result T
		if err = json.Unmarshal(assetState.Value, &result); err != nil {
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
// Delete (only for testing)
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
	var emptyObject T

	compositeKey, err := ctx.GetStub().CreateCompositeKey(emptyObject.Type(), []string{key})
	if err != nil {
		return &CompositeKeyCreationError{err.Error(), key, emptyObject.Type()}
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
// Custom Methods
// =============================================================================

// Casts a vote for a ballot.
// This function will assert that the ballot has been assigned to voterID and has a matching candidate with candidateID.
// This function will return an error if the vote has already been cast.
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

	if err = ballot.Vote(candidateID); err != nil {
		return err
	}

	return updateAsset(ctx, ballot.Asset.ID, ballot)
}

// Helper function to sync the election and candidates. Duplicates are aptly handled.
func (s *SmartContract) SyncElectionAndCandidates(ctx contractapi.TransactionContextInterface, electionID string) error {
	election, err := queryAsset[Election](ctx, electionID)
	if err != nil {
		return err
	}

	allCandidates, err := queryAssetsByType[Candidate](ctx)
	if err != nil {
		return err
	}

	// Update election with candidates assigned to it
	// We clear the candidate slice in the election to prevent duplicates
	election.Candidates = election.Candidates[:0]
	for i := range allCandidates {
		if allCandidates[i].ElectionID != electionID {
			continue
		}

		election.Candidates = append(election.Candidates, allCandidates[i].Asset.ID)
	}

	if err = updateAsset(ctx, election.Asset.ID, election); err != nil {
		return err
	}

	return nil
}

// =============================================================================
// Performance Testing
// =============================================================================

func (s *SmartContract) SetupPerformanceTestElection(ctx contractapi.TransactionContextInterface, voterID string, numBallots int, numCandidates int) (string, error) {

	// Setup Election with no candidates
	electionID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	election := Election{
		Asset:     Asset{ID: "e-" + electionID.String()},
		EndTime:   "2024-03-24 23:59:59",
		StartTime: "2024-03-23 00:00:00",
	}

	// Create candidates
	candidates := []Candidate{}

	for i := 0; i < numCandidates; i++ {
		candidateID, err := uuid.NewV7()
		if err != nil {
			return "", err
		}

		candidate := Candidate{
			Asset:      Asset{ID: "c-" + candidateID.String()},
			ElectionID: election.Asset.ID,
			Name:       fmt.Sprintf("performanceTestCandidate%d", i),
			PublicKey:  "eyJOIjoxMDA5ODc4ODk1NjcwMjc2NjY5OTQzMzg2Njk3MjUzNDA1ODMxNjE2NDY3MjMwODgyNDg5MjIyMzI4ODc4NjI3NDE4Nzg2MDExODQ3MDcsIk5TcXVhcmUiOjEwMTk4NTUzODM5MjAyMTc1NTEwMjI2ODQ5NTUwMzAzMzM3ODAzNTY2MTUwNjE4NTE0NDY0NTM1NTk1NDEwNjE0MzQ3NzIxNjU5NjgyNjY2ODE1NTExMTIzMzgyNzcxMzczMjE5NzA4OTYzMzQwOTcxODY2NzYxNjM4NzE4NDA5MDA2MDE5NTQ4NjQ1NTQyNTQzOTMwNjc1ODQ5LCJHIjoxMDA5ODc4ODk1NjcwMjc2NjY5OTQzMzg2Njk3MjUzNDA1ODMxNjE2NDY3MjMwODgyNDg5MjIyMzI4ODc4NjI3NDE4Nzg2MDExODQ3MDgsIkxlbmd0aCI6MTI4fQ==",
		}

		election.Candidates = append(election.Candidates, candidate.Asset.ID)
		candidates = append(candidates, candidate)

		if err = candidate.Init(); err != nil {
			return "", err
		}
		if err = createAsset(ctx, candidate.Asset.ID, candidate); err != nil {
			return "", err
		}

	}

	// Create election
	if err = createAsset(ctx, election.Asset.ID, election); err != nil {
		return "", err
	}

	// Create 100 ballots
	for i := 0; i < numBallots; i++ {
		ballotID, err := uuid.NewV7()
		if err != nil {
			return "", err
		}

		ballot := Ballot{
			Asset:      Asset{ID: "b-" + ballotID.String()},
			Candidates: candidates,
			ElectionID: election.Asset.ID,
			VoterID:    voterID,
		}

		for i := range ballot.Candidates {
			if err = ballot.Candidates[i].Init(); err != nil {
				return "", err
			}
		}

		if err = createAsset(ctx, ballot.Asset.ID, ballot); err != nil {
			return "", err
		}
	}

	return election.Asset.ID, nil
}

func (s *SmartContract) CastVotesForPerformanceTestElection(ctx contractapi.TransactionContextInterface, electionID string) error {
	election, err := queryAsset[Election](ctx, electionID)
	if err != nil {
		return err
	}

	var emptyBallot Ballot

	// Get all ballots from this election
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(emptyBallot.Type(), []string{})
	if err != nil {
		return err
	}
	defer resultIterator.Close()

	for resultIterator.HasNext() {
		assetState, err := resultIterator.Next()
		if err != nil {
			return err
		}

		var ballot Ballot
		if err = json.Unmarshal(assetState.Value, &ballot); err != nil {
			return err
		}

		// Ensure it matches this election
		if ballot.ElectionID == electionID {
			// RNG the candidateID
			numCandidates := int64(len(election.Candidates))

			index, err := rand.Int(rand.Reader, big.NewInt(numCandidates))
			if err != nil {
				return err
			}

			candidateID := election.Candidates[index.Int64()]

			// Cast Vote
			if err = ballot.Vote(candidateID); err != nil {
				return err
			}

			// Update ballot
			if err = updateAsset(ctx, ballot.Asset.ID, ballot); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SmartContract) CleanUpPerformanceTestElection(ctx contractapi.TransactionContextInterface, electionID string) error {
	// Get all ballot IDs
	var emptyBallot Ballot

	// Get all ballots from this election
	ballotIDsToDelete := []string{}
	resultIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(emptyBallot.Type(), []string{})
	if err != nil {
		return err
	}
	defer resultIterator.Close()

	for resultIterator.HasNext() {
		assetState, err := resultIterator.Next()
		if err != nil {
			return err
		}

		var ballot Ballot
		if err = json.Unmarshal(assetState.Value, &ballot); err != nil {
			return err
		}

		if ballot.ElectionID == electionID {
			ballotIDsToDelete = append(ballotIDsToDelete, ballot.Asset.ID)
		}
	}

	// Get candidate IDs
	election, err := queryAsset[Election](ctx, electionID)
	if err != nil {
		return err
	}

	// Delete ballots
	for _, id := range ballotIDsToDelete {
		if err = deleteAsset[Ballot](ctx, id); err != nil {
			return err
		}
	}

	// Delete candidates
	for _, id := range election.Candidates {
		if err = deleteAsset[Candidate](ctx, id); err != nil {
			return err
		}
	}

	// Delete election
	if err = deleteAsset[Election](ctx, electionID); err != nil {
		return err
	}

	return nil
}
