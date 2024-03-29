package common

// =============================================================================
// DynamoDB Items
// =============================================================================

type VoterCredentials struct {
	NRIC       string `json:"nric"`
	ElectionID string `json:"electionID"`
	VoterID    string `json:"voterID"`
	BallotID   string `json:"ballotID"`
}
