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

	t.Run("Successfully create ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		_, mockBallotData := MockBallot()
		_, mockElectionData := MockElection()

		mockStub.On("GetState", "b-0").Return(nil, nil)
		mockStub.On("GetState", "e-0").Return(mockElectionData, nil)
		mockStub.On("PutState", "b-0", mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.NoError(t, err)
	})

	t.Run("Fail to create existing ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		_, mockBallotData := MockBallot()
		_, mockElectionData := MockElection()

		mockStub.On("GetState", "e-0").Return(mockElectionData, nil)
		mockStub.On("GetState", "b-0").Return(mockBallotData, nil)

		// Test
		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, "ballot: b-0 already created")
	})

	t.Run("Fail to create invalid ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		// Modify ballot for fail case
		mockBallot, _ := MockBallot()
		mockBallot.BallotID = ""
		mockBallotData, err := json.Marshal(mockBallot)
		if err != nil {
			t.Error(err)
		}

		// Test
		// We don't need to mock GetState for the ballotID since it will fail before reaching createAsset
		expectedError := fmt.Sprintf("%s is invalid! %s", reflect.TypeOf(*mockBallot).String(), "missing BallotID")

		err = smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, expectedError)
	})

	t.Run("Fail to create ballot for non-existent election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		_, mockBallotData := MockBallot()
		mockStub.On("GetState", "e-0").Return(nil, nil)

		// Test
		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, "cannot read world state with key e-0")
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

func TestQuery(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("successfully query ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, mockBallotData := MockBallot()

		mockStub.On("GetState", "b-0").Return(mockBallotData, nil)

		// Test
		ballot, err := smartContract.QueryBallot(mockCtx, mockBallot.BallotID)
		if err != nil {
			t.Error(err)
		}
		require.Equal(t, *mockBallot, ballot)
	})

	t.Run("fail to query non-existent ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, _ := MockBallot()

		mockStub.On("GetState", "b-0").Return(nil, nil)

		// Test
		_, err := smartContract.QueryBallot(mockCtx, mockBallot.BallotID)
		require.EqualError(t, err, "cannot read world state with key b-0")
	})
}

// =============================================================================
// Update Tests
// =============================================================================

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

func MockCandidate() (*chaincode.Candidate, []byte) {
	mock := chaincode.Candidate{
		CandidateID: "c-0",
		Count:       0,
		ElectionID:  "e-0",
		Name:        "mockCandidate",
	}

	mockData, err := json.Marshal(mock)
	if err != nil {
		log.Fatal(err)
	}

	return &mock, mockData
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

func MockVoter() (*chaincode.Voter, []byte) {
	mock := chaincode.Voter{
		VoterID:  "v-0",
		BallotID: "",
	}

	mockData, err := json.Marshal(mock)
	if err != nil {
		log.Fatal(err)
	}

	return &mock, mockData
}
