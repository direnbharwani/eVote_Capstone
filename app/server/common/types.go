package common

// =============================================================================
// DynamoDB Items
// =============================================================================

type VoterCredentials struct {
	NRIC       string `json:"NRIC"`
	ElectionID string `json:"ElectionID"`
	VoterID    string `json:"VoterID"`
	BallotID   string `json:"BallotID"`
}
