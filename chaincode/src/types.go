package chaincode

import (
	"log"
	"time"
)

// =============================================================================
// BioData
// =============================================================================

// BioData for a citizen
type BioData struct {
	Birthday string `json:"Birthday"`
	DID      string `json:"DID"`
	Name     string `json:"Name"`
	NRIC     string `json:"NRIC"`
	Sex      string `json:"Sex"`
}

func (b *BioData) GetAge() int {
	birthday, err := time.Parse(time.DateOnly, b.Birthday)
	if err != nil {
		log.Printf("Failed to parse Birthday for %s: %s", b.DID, err.Error())
	}

	now := time.Now()
	age := now.Year() - birthday.Year()

	// Modify age if birthday has not been reached
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}
	return age
}

func (b *BioData) EligibleToVote() bool {
	return b.GetAge() < 21
}

// =============================================================================
// Election
// =============================================================================

type Election struct {
	Candidates []Candidate `json:"Candidates"`
	ElectionID string      `json:"ElectionID"`
	EndTime    string      `json:"EndTime"`
	StartTime  string      `json:"StartTime"`
}

func (e *Election) IsActive() bool {
	now := time.Now()

	start, err := time.Parse(time.DateTime, e.StartTime)
	if err != nil {
		log.Fatal(err)
	}

	end, err := time.Parse(time.DateTime, e.EndTime)
	if err != nil {
		log.Fatal(err)
	}

	return (now.After(start) && now.Before(end))
}

// =============================================================================
// Candidate
// =============================================================================

// Defines a electoral candidate
type Candidate struct {
	BioData     BioData `json:"BioData"`
	CandidateID string  `json:"CandidateID"`
	Count       uint64  `json:"Count"`
}

// =============================================================================
// Voter
// =============================================================================

// Defines a Voter that is created with a ballot
// Asset ID for Voters are prefixed with v-
type Voter struct {
	BioData  BioData `json:"BioData"`
	VoterID  string  `json:"VoterID"`
	BallotID string  `json:"BallotID"`
}

// =============================================================================
// Ballot
// =============================================================================

// Defines a Ballot that is assigned to a voter
// Asset ID for Ballots are prefixed with b-
type Ballot struct {
	BallotID   string      `json:"BallotID"`
	Candidates []Candidate `json:"Candidates"`
	ElectionID string      `json:"ElectionID"`
	VoterID    string      `json:"VoterID"`
	Voted      bool        `json:"Voted"`
}
