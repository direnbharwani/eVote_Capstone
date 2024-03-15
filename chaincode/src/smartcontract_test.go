package chaincode_test

import (
	"testing"

	chaincode "github.com/direnbharwani/eVote_Capstone/src"
	mocks "github.com/direnbharwani/eVote_Capstone/src/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Mock Objects
// =============================================================================

func MockCandidate() *chaincode.Candidate {
	mockCandidateBioData := chaincode.BioData{
		Birthday: "01-01-2024",
		DID:      "1",
		Name:     "mockCandidate",
		NRIC:     "12345",
		Sex:      "M",
	}

	return &chaincode.Candidate{
		BioData:     mockCandidateBioData,
		CandidateID: "c-0",
		Count:       0,
	}
}

func MockElection() *chaincode.Election {
	return &chaincode.Election{
		Candidates: []chaincode.Candidate{*MockCandidate()},
		ElectionID: "e-0",
		Name:       "mockElection",
		EndTime:    "2024-01-01 23:59:59",
		StartTime:  "2024-01-01 00:00:00",
	}
}

func MockVoter() *chaincode.Voter {
	mockVoterBioData := chaincode.BioData{
		Birthday: "01-01-2024",
		DID:      "1",
		Name:     "mockVoter",
		NRIC:     "12345",
		Sex:      "M",
	}

	return &chaincode.Voter{
		BioData:  mockVoterBioData,
		VoterID:  "v-0",
		BallotID: "",
	}
}

// =============================================================================
// Creation Tests
// =============================================================================

func TestCreateCandidate(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	// Mocks for Interfaces
	mockStub := &mocks.ChaincodeStubInterface{}
	mockCtx := &mocks.TransactionContextInterface{}

	mockStub.On("GetState", "c-0").Return(nil, nil)
	mockStub.On("PutState", "c-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

	mockCtx.On("GetStub").Return(mockStub)

	// Test
	err := smartContract.CreateCandidate(mockCtx, *MockCandidate())
	require.NoError(t, err)
}

func TestCreateElection(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	// Mocks for Interfaces
	mockStub := &mocks.ChaincodeStubInterface{}
	mockCtx := &mocks.TransactionContextInterface{}

	mockCtx.On("GetStub").Return(mockStub)

	mockStub.On("GetState", "e-0").Return(nil, nil)
	mockStub.On("PutState", "e-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

	// Test
	err := smartContract.CreateElection(mockCtx, *MockElection())
	require.NoError(t, err)
}

func TestCreateVoter(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	// Mocks for Interfaces
	mockStub := &mocks.ChaincodeStubInterface{}
	mockCtx := &mocks.TransactionContextInterface{}

	mockCtx.On("GetStub").Return(mockStub)

	mockStub.On("GetState", "v-0").Return(nil, nil)
	mockStub.On("PutState", "v-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

	// Test
	err := smartContract.CreateVoter(mockCtx, *MockVoter())
	require.NoError(t, err)
}
