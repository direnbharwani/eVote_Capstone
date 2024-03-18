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

func TestCreateCandidate(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully create candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockCandidate, mockCandidateData := MockCandidate()

		mockStub.On("CreateCompositeKey", mockCandidate.Type(), []string{mockCandidate.Asset.ID}).Return(mockCandidate.Asset.ID, nil)
		mockStub.On("GetState", mockCandidate.Asset.ID).Return(nil, nil)
		mockStub.On("PutState", mockCandidate.Asset.ID, mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.CreateCandidate(mockCtx, string(mockCandidateData))
		require.NoError(t, err)
	})

	t.Run("Fail to create existing candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockCandidate, mockCandidateData := MockCandidate()

		mockStub.On("CreateCompositeKey", mockCandidate.Type(), []string{mockCandidate.Asset.ID}).Return(mockCandidate.Asset.ID, nil)
		mockStub.On("GetState", mockCandidate.Asset.ID).Return(mockCandidateData, nil)

		// Test
		expectedError := fmt.Sprintf("%s: %s already created", mockCandidate.Type(), mockCandidate.Asset.ID)

		err := smartContract.CreateCandidate(mockCtx, string(mockCandidateData))
		require.EqualError(t, err, expectedError)
	})

	t.Run("Fail to create invalid candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		// Modify candidate for fail case
		mockCandidate, _ := MockCandidate()
		mockCandidate.Asset.ID = ""
		mockCandidateData, err := json.Marshal(mockCandidate)
		if err != nil {
			t.Error(err)
		}

		// Test
		expectedError := fmt.Sprintf("%s is invalid! %s", mockCandidate.Type(), "missing ID")

		err = smartContract.CreateCandidate(mockCtx, string(mockCandidateData))
		require.EqualError(t, err, expectedError)
	})
}

func TestCreateElection(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully create election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(nil, nil)
		mockStub.On("PutState", mockElection.Asset.ID, mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.CreateElection(mockCtx, string(mockElectionData))
		require.NoError(t, err)
	})

	t.Run("Fail to create existing election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(mockElectionData, nil)

		// Test
		expectedError := fmt.Sprintf("%s: %s already created", mockElection.Type(), mockElection.Asset.ID)

		err := smartContract.CreateElection(mockCtx, string(mockElectionData))
		require.EqualError(t, err, expectedError)
	})

	t.Run("Fail to create invalid candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		// Modify election for fail case
		mockElection, _ := MockElection()
		mockElection.StartTime = "error"
		mockElectionData, err := json.Marshal(mockElection)
		if err != nil {
			t.Error(err)
		}

		// Test
		expectedError := fmt.Sprintf("%s is invalid! %s", mockElection.Type(), "parsing time \"error\" as \"2006-01-02 15:04:05\": cannot parse \"error\" as \"2006\"")

		err = smartContract.CreateElection(mockCtx, string(mockElectionData))
		require.EqualError(t, err, expectedError)
	})
}

// =============================================================================
// Query Tests
// =============================================================================

func TestQueryBallot(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully query ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, mockBallotData := MockBallot()

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(mockBallotData, nil)

		// Test
		result, err := smartContract.QueryBallot(mockCtx, mockBallot.Asset.ID)
		if err != nil {
			t.Error(err)
		}
		require.Equal(t, *mockBallot, result)
	})

	t.Run("Fail to query non-existent ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, _ := MockBallot()

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockBallot.Asset.ID)

		_, err := smartContract.QueryBallot(mockCtx, mockBallot.Asset.ID)
		require.EqualError(t, err, expectedError)
	})
}

func TestQueryCandidate(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully query candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockCandidate, mockCandidateData := MockCandidate()

		mockStub.On("CreateCompositeKey", mockCandidate.Type(), []string{mockCandidate.Asset.ID}).Return(mockCandidate.Asset.ID, nil)
		mockStub.On("GetState", mockCandidate.Asset.ID).Return(mockCandidateData, nil)

		// Test
		result, err := smartContract.QueryCandidate(mockCtx, mockCandidate.Asset.ID)
		if err != nil {
			t.Error(err)
		}
		require.Equal(t, *mockCandidate, result)
	})

	t.Run("Fail to query non-existent candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockCandidate, _ := MockCandidate()

		mockStub.On("CreateCompositeKey", mockCandidate.Type(), []string{mockCandidate.Asset.ID}).Return(mockCandidate.Asset.ID, nil)
		mockStub.On("GetState", mockCandidate.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockCandidate.Asset.ID)

		_, err := smartContract.QueryCandidate(mockCtx, mockCandidate.Asset.ID)
		require.EqualError(t, err, expectedError)
	})
}

func TestQueryElection(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully query election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(mockElectionData, nil)

		// Test
		result, err := smartContract.QueryElection(mockCtx, mockElection.Asset.ID)
		if err != nil {
			t.Error(err)
		}
		require.Equal(t, *mockElection, result)
	})

	t.Run("Fail to query non-existent election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockElection, _ := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockElection.Asset.ID)

		_, err := smartContract.QueryElection(mockCtx, mockElection.Asset.ID)
		require.EqualError(t, err, expectedError)
	})
}

// =============================================================================
// Update Tests
// =============================================================================

func TestUpdateBallot(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully update ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, mockBallotData := MockBallot()

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(mockBallotData, nil)
		mockStub.On("PutState", mockBallot.Asset.ID, mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.UpdateBallot(mockCtx, string(mockBallotData))
		require.NoError(t, err)
	})

	t.Run("Fail to update non-existent ballot", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockBallot, mockBallotData := MockBallot()

		mockStub.On("CreateCompositeKey", mockBallot.Type(), []string{mockBallot.Asset.ID}).Return(mockBallot.Asset.ID, nil)
		mockStub.On("GetState", mockBallot.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockBallot.Asset.ID)

		err := smartContract.UpdateBallot(mockCtx, string(mockBallotData))
		require.EqualError(t, err, expectedError)
	})
}

func TestUpdateCandidate(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully update candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockCandidate, mockCandidateData := MockCandidate()

		mockStub.On("CreateCompositeKey", mockCandidate.Type(), []string{mockCandidate.Asset.ID}).Return(mockCandidate.Asset.ID, nil)
		mockStub.On("GetState", mockCandidate.Asset.ID).Return(mockCandidateData, nil)
		mockStub.On("PutState", mockCandidate.Asset.ID, mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.UpdateCandidate(mockCtx, string(mockCandidateData))
		require.NoError(t, err)
	})

	t.Run("Fail to update non-existent candidate", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockCandidate, mockCandidateData := MockCandidate()

		mockStub.On("CreateCompositeKey", mockCandidate.Type(), []string{mockCandidate.Asset.ID}).Return(mockCandidate.Asset.ID, nil)
		mockStub.On("GetState", mockCandidate.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockCandidate.Asset.ID)

		err := smartContract.UpdateCandidate(mockCtx, string(mockCandidateData))
		require.EqualError(t, err, expectedError)
	})
}

func TestUpdateElection(t *testing.T) {
	smartContract := chaincode.SmartContract{}

	t.Run("Successfully update Election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(mockElectionData, nil)
		mockStub.On("PutState", mockElection.Asset.ID, mock.AnythingOfType("[]uint8")).Return(nil, nil)

		// Test
		err := smartContract.UpdateElection(mockCtx, string(mockElectionData))
		require.NoError(t, err)
	})

	t.Run("Fail to update non-existent Election", func(t *testing.T) {
		// Mocks
		mockStub := &mocks.ChaincodeStubInterface{}
		mockCtx := &mocks.TransactionContextInterface{}

		mockCtx.On("GetStub").Return(mockStub)

		mockElection, mockElectionData := MockElection()

		mockStub.On("CreateCompositeKey", mockElection.Type(), []string{mockElection.Asset.ID}).Return(mockElection.Asset.ID, nil)
		mockStub.On("GetState", mockElection.Asset.ID).Return(nil, nil)

		// Test
		expectedError := fmt.Sprintf("cannot read world state with key %s", mockElection.Asset.ID)

		err := smartContract.UpdateElection(mockCtx, string(mockElectionData))
		require.EqualError(t, err, expectedError)
	})
}

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
