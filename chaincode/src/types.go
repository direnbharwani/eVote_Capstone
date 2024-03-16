package chaincode

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

// ITYPES is a union set type constraint
// that enforces only allowable types are passed to smart contract methods.
type ITYPES interface {
	Ballot | Candidate | Election | Voter

	// Methods
	// Type() string
	IsValid() error
}

// =============================================================================
// Errors
// =============================================================================

type ObjectValidationError struct {
	Message string
	Type    string
}

func (e *ObjectValidationError) Error() string {
	return fmt.Sprintf("%s is invalid! %s", e.Type, e.Message)
}

// =============================================================================
// Election
// =============================================================================

// Defines an election
// Asset ID for Elections are prefixed with e-
type Election struct {
	Candidates []Candidate `json:"Candidates"`
	ElectionID string      `json:"ElectionID"`
	EndTime    string      `json:"EndTime"`
	Name       string      `json:"Name"`
	StartTime  string      `json:"StartTime"`
}

func (e Election) IsValid() error {
	objectType := reflect.TypeOf(e).String()

	if e.ElectionID == "" {
		return &ObjectValidationError{"missing ElectionID", objectType}
	}

	if _, err := time.Parse(time.DateTime, e.StartTime); err != nil {
		return &ObjectValidationError{err.Error(), objectType}
	}

	if _, err := time.Parse(time.DateTime, e.EndTime); err != nil {
		return &ObjectValidationError{err.Error(), objectType}
	}

	return nil
}

func (e Election) IsActive() bool {
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
// Asset ID for Candidates are prefixed with c-
type Candidate struct {
	CandidateID string `json:"CandidateID"`
	Count       uint64 `json:"Count"`
	ElectionID  string `json:"ElectionID"`
	Name        string `json:"Name"`
}

func (c Candidate) IsValid() error {
	objectType := reflect.TypeOf(c).String()

	if c.CandidateID == "" {
		return &ObjectValidationError{"missing CandidateID", objectType}
	}

	if c.ElectionID == "" {
		return &ObjectValidationError{"missing ElectionID", objectType}
	}

	return nil
}

// =============================================================================
// Voter
// =============================================================================

// Defines a Voter that is created with a ballot
// Asset ID for Voters are prefixed with v-
type Voter struct {
	VoterID  string `json:"VoterID"`
	BallotID string `json:"BallotID"`
}

func (v Voter) IsValid() error {
	objectType := reflect.TypeOf(v).String()

	if v.VoterID == "" {
		return &ObjectValidationError{"missing VoterID", objectType}
	}

	return nil
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

// Checks if BallotID, VoterID & ElectionID are not empty strings
func (b Ballot) IsValid() error {
	objectType := reflect.TypeOf(b).String()

	if b.BallotID == "" {
		return &ObjectValidationError{"missing BallotID", objectType}
	}

	if b.VoterID == "" {
		return &ObjectValidationError{"missing VoterID", objectType}
	}

	if b.ElectionID == "" {
		return &ObjectValidationError{"missing ElectionID", objectType}
	}

	return nil
}
