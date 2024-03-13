package chaincode

import (
	"encoding/json"
	"fmt"

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
// Registration
// =============================================================================

func (s *SmartContract) RegisterCandidate(ctx contractapi.TransactionContextInterface, candidate Candidate) error {
	return createAsset[Candidate](ctx, candidate.CandidateID, candidate, "candidate")
}

func (s *SmartContract) RegisterElection(ctx contractapi.TransactionContextInterface, election Election) error {
	return createAsset[Election](ctx, election.ElectionID, election, "election")
}

func (s *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface, voter Voter) error {
	return createAsset[Voter](ctx, voter.VoterID, voter, "voter")
}

func (s *SmartContract) CreateBallot(ctx contractapi.TransactionContextInterface, voterID string, electionID string) error {
	// Check if voter has an existing ballot
	voterState, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}

	if voterState == nil {
		return fmt.Errorf("voter %s does not exist", voterID)
	}

	voter := Voter{}
	err = json.Unmarshal(voterState, &voter)
	if err != nil {
		return err
	}

	if voter.BallotID != "" { // existing ballot
		return fmt.Errorf("voter %s has an existing ballot %s", voterID, voter.BallotID)
	}

	// TODO: Create ballot

	return nil
}

func createAsset[T ITYPES](ctx contractapi.TransactionContextInterface, key string, body T, assetType string) error {
	objectState, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}
	if objectState != nil {
		return fmt.Errorf("%s: %s already created", assetType, key)
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(key, bytes)
	if err != nil {
		return fmt.Errorf("unable to interact with world state: %v", err)
	}

	return nil
}

// =============================================================================
// Query
// =============================================================================

// =============================================================================
// Update
// =============================================================================
