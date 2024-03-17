package chaincode_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	chaincode "github.com/direnbharwani/eVote_Capstone/src"
	mocks "github.com/direnbharwani/eVote_Capstone/src/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Creation Tests
// =============================================================================

func TestCreateBallot(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully create ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, mockBallotData := MockBallot()
		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)

		mockStub.On("GetState", mockElection.Asset.ID).Return(mockElectionData, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(nil, nil)
		mockStub.On("PutState", mockBallot.Asset.ID, mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.NoError(t, err)
	})

	t.Run("Fail to create existing ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, mockBallotData := MockBallot()
		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)

		mockStub.On("GetState", mockElection.Asset.ID).Return(mockElectionData, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(mockBallotData, nil)

		// Test
		expectedError := fmt.Sprintf("%s: %s already created", mockBallot.Type(), mockBallot.Asset.ID)

		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, expectedError)
	})

	t.Run("Fail to create invalid ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		// Modify ballot for fail case
		mockBallot, _ := MockBallot()
		mockBallot.Asset.ID = ""
		mockBallotData, err := json.Marshal(mockBallot)
		if err != nil {
			t.Error(err)
		}

		// Test
		// We don't need to mock GetState for the ballotID since it will fail before reaching createAsset
		expectedError := fmt.Sprintf("%s is invalid! %s", mockBallot.Type(), "missing ID")

		err = smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, expectedError)
	})

	t.Run("Fail to create ballot for non-existent election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		_, mockBallotData := MockBallot()
		mockElection, _ := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockElection.Asset.ID)

		err := smartContract.CreateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, expectedError)
	})
}

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

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(mockBallotData, nil)

		// Test
		ballot, err := smartContract.QueryBallot(mockCtx, mockBallot.Asset.ID)
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

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(nil, nil)

		// Test
		_, err := smartContract.QueryBallot(mockCtx, mockBallot.Asset.ID)
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
	id := chaincode.Asset{"b-0"}

	mock := chaincode.Ballot{
		Asset:      id,
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
	id := chaincode.Asset{"c-0"}

	mock := chaincode.Candidate{
		Asset:      id,
		Count:      0,
		ElectionID: "e-0",
		Name:       "mockCandidate",
	}

	mockData, err := json.Marshal(mock)
	if err != nil {
		log.Fatal(err)
	}

	return &mock, mockData
}

func MockElection() (*chaincode.Election, []byte) {
	id := chaincode.Asset{"e-0"}

	mock := chaincode.Election{
		Asset:      id,
		Candidates: []chaincode.Candidate{},
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
	id := chaincode.Asset{"v-0"}

	mock := chaincode.Voter{
		Asset:    id,
		BallotID: "",
	}

	mockData, err := json.Marshal(mock)
	if err != nil {
		log.Fatal(err)
	}

	return &mock, mockData
}
