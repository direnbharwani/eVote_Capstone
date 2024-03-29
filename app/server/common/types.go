package common

// =============================================================================
// DynamoDB Items
// =============================================================================

type VoterCredentials struct {
	NRIC       string `dynamodbav:"nric"`
	ElectionID string `dynamodbav:"electionID"`
	VoterID    string `dynamodbav:"voterID"`
	BallotID   string `dynamodbav:"ballotID"`
}
