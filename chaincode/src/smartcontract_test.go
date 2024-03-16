package chaincode_test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	chaincode "github.com/direnbharwani/eVote_Capstone/src"
	mocks "github.com/direnbharwani/eVote_Capstone/src/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Creation Tests
// =============================================================================

func TestCreation(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	// Mocks
	mockStub := &mocks.ChaincodeStubInterface{}
	mockCtx := &mocks.TransactionContextInterface{}

	mockCtx.On("GetStub").Return(mockStub)

	t.Run("successfully create ballot", func(t *testing.T) {
		_, mockBallotData := MockBallot()
		_, mockElectionData := MockElection()

		mockStub.On("GetState", "b-0").Return(nil, nil)
		mockStub.On("GetState", "e-0").Return(mockElectionData, nil)
		mockStub.On("PutState", "b-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.NoError(t, err)
	})

	t.Run("failed to create invalid ballot", func(t *testing.T) {
		// Modify ballot for fail case
		mockBallot, _ := MockBallot()
		mockBallot.BallotID = ""
		mockBallotData, err := json.Marshal(mockBallot)
		if err != nil {
			t.Error(err)
		}

		// We don't need to mock GetState for the ballotID since it will fail before reaching createAsset
		expectedError := fmt.Sprintf("%s is invalid! %s", reflect.TypeOf(*mockBallot).String(), "missing BallotID")

		err = smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, expectedError)
	})
}

// func TestCreateCandidate(t *testing.T) {
// 	smartContract := chaincode.SmartContract{}

// 	// Mocks
// 	mockStub := &mocks.ChaincodeStubInterface{}
// 	mockCtx := &mocks.TransactionContextInterface{}

// 	mockCtx.On("GetStub").Return(mockStub)

// 	mockStub.On("GetState", "c-0").Return(nil, nil)
// 	mockStub.On("PutState", "c-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

// 	// Test
// 	err := smartContract.CreateCandidate(mockCtx, *MockCandidate())
// 	require.NoError(t, err)
// }

// func TestCreateElection(t *testing.T) {
// 	smartContract := chaincode.SmartContract{}

// 	// Mocks
// 	mockStub := &mocks.ChaincodeStubInterface{}
// 	mockCtx := &mocks.TransactionContextInterface{}

// 	mockCtx.On("GetStub").Return(mockStub)

// 	mockStub.On("GetState", "e-0").Return(nil, nil)
// 	mockStub.On("PutState", "e-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

// 	// Test
// 	err := smartContract.CreateElection(mockCtx, *MockElection())
// 	require.NoError(t, err)
// }

// func TestCreateVoter(t *testing.T) {
// 	smartContract := chaincode.SmartContract{}

// 	// Mocks
// 	mockStub := &mocks.ChaincodeStubInterface{}
// 	mockCtx := &mocks.TransactionContextInterface{}

// 	mockCtx.On("GetStub").Return(mockStub)

// 	mockStub.On("GetState", "v-0").Return(nil, nil)
// 	mockStub.On("PutState", "v-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

// 	// Test
// 	err := smartContract.CreateVoter(mockCtx, *MockVoter())
// 	require.NoError(t, err)
// }

// =============================================================================
// Query Tests
// =============================================================================

// func TestQueryBallot(t *testing.T) {
// 	smartContract := chaincode.SmartContract{}

// 	// Mocks for Interfaces
// 	mockStub := &mocks.ChaincodeStubInterface{}
// 	mockCtx := &mocks.TransactionContextInterface{}

// 	mockCtx.On("GetStub").Return(mockStub)

// 	mockBallot := *MockBallot()
// 	mockBallotData, err := json.Marshal(mockBallot)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	mockStub.On("GetState", mockBallot.BallotID).Return(mockBallotData, nil)

// 	// Test
// 	want := mockBallot

// 	got, err := smartContract.QueryBallot(mockCtx, mockBallot.BallotID)
// 	require.NoError(t, err)

// 	assertBallotEquality(t, got, want)
// }

// =============================================================================
// Update Tests
// =============================================================================

// =============================================================================
// Assertion Helpers
// =============================================================================

// Defined due to array of custom objects
func assertBallotEquality(t testing.TB, lhs chaincode.Ballot, rhs chaincode.Ballot) {
	if lhs.BallotID != rhs.BallotID {
		t.Errorf("got %v want %v", lhs, rhs)
	}

	if lhs.ElectionID != rhs.ElectionID {
		t.Errorf("got %v want %v", lhs, rhs)
	}

	if lhs.VoterID != rhs.VoterID {
		t.Errorf("got %v want %v", lhs, rhs)
	}

	if lhs.Voted != rhs.Voted {
		t.Errorf("got %v want %v", lhs, rhs)
	}

	// Compare all candidates
	if len(lhs.Candidates) != len(rhs.Candidates) {
		t.Errorf("got %v want %v", lhs, rhs)
	}

	for i, candidate := range lhs.Candidates {
		if candidate != rhs.Candidates[i] {
			t.Errorf("got %v want %v", lhs, rhs)
		}
	}
}

func assertElectionEquality(t testing.TB, lhs chaincode.Election, rhs chaincode.Election) {
}

// =============================================================================
// Mock Objects
// =============================================================================

func MockBallot() (*chaincode.Ballot, []byte) {
	mock := chaincode.Ballot{
		BallotID:   "b-0",
		Candidates: []chaincode.Candidate{},
		ElectionID: "e-0",
		VoterID:    "v-0",
		Voted:      false,
	}

	mockData, err := json.Marshal(mock)
	if err != nil {
		log.Fatal(err)
	}

	return &mock, mockData
}

func MockCandidate() *chaincode.Candidate {
	return &chaincode.Candidate{
		CandidateID: "c-0",
		Count:       0,
		Name:        "mockCandidate",
	}
}

func MockElection() (*chaincode.Election, []byte) {
	mock := chaincode.Election{
		Candidates: []chaincode.Candidate{},
		ElectionID: "e-0",
		Name:       "mockElection",
		EndTime:    "2024-01-01 23:59:59",
		StartTime:  "2024-01-01 00:00:00",
	}

	mockData, err := json.Marshal(mock)
	if err != nil {
		log.Fatal(err)
	}

	return &mock, mockData
}

func MockVoter() *chaincode.Voter {
	return &chaincode.Voter{
		VoterID:  "v-0",
		BallotID: "",
	}
}
