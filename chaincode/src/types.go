package chaincode

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

// ITYPES is a union set type constraint
// that enforces only allowable types are passed to smart contract methods.
type ITYPES interface {
	Ballot | Candidate | Election | Voter

	Type() string
	IsValid() error
}

type Asset struct {
	ID string `json:"ID"`
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
	Asset      Asset       `json:"Asset"`
	Candidates []Candidate `json:"Candidates"`
	EndTime    string      `json:"EndTime"`
	Name       string      `json:"Name"`
	StartTime  string      `json:"StartTime"`
}

func (e Election) Type() string {
	return reflect.TypeOf(e).String()
}

func (e Election) IsValid() error {
	objectType := reflect.TypeOf(e).String()

	if e.Asset.ID == "" {
		return &ObjectValidationError{"missing ID", objectType}
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
	Asset      Asset  `json:"Asset"`
	Count      uint64 `json:"Count"`
	ElectionID string `json:"ElectionID"`
	Name       string `json:"Name"`
}

func (c Candidate) Type() string {
	return reflect.TypeOf(c).String()
}

func (c Candidate) IsValid() error {
	objectType := reflect.TypeOf(c).String()

	if c.Asset.ID == "" {
		return &ObjectValidationError{"missing ID", objectType}
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
	Asset    Asset  `json:"Asset"`
	BallotID string `json:"BallotID"`
}

func (v Voter) Type() string {
	return reflect.TypeOf(v).String()
}

func (v Voter) IsValid() error {
	objectType := reflect.TypeOf(v).String()

	if v.Asset.ID == "" {
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
	Asset      Asset       `json:"Asset"`
	Candidates []Candidate `json:"Candidates"`
	ElectionID string      `json:"ElectionID"`
	VoterID    string      `json:"VoterID"`
	Voted      bool        `json:"Voted"`
}

func (b Ballot) Type() string {
	return reflect.TypeOf(b).String()
}

// Checks if BallotID, VoterID & ElectionID are not empty strings
func (b Ballot) IsValid() error {
	objectType := reflect.TypeOf(b).String()

	if b.Asset.ID == "" {
		return &ObjectValidationError{"missing ID", objectType}
	}

	if b.VoterID == "" {
		return &ObjectValidationError{"missing VoterID", objectType}
	}

	if b.ElectionID == "" {
		return &ObjectValidationError{"missing ElectionID", objectType}
	}

	return nil
}

func (b *Ballot) Vote(candidateID string) error {
	candidateFound := false
	for i, c := range b.Candidates {
		if c.Asset.ID == candidateID {
			candidateFound = true

			b.Candidates[i].Count++
			b.Voted = true

			break
		}
	}

	if !candidateFound {
		errorMessage := fmt.Sprintf("candidate %s is not found in ballot %s!", candidateID, b.Asset.ID)
		return errors.New(errorMessage)
	}

	return nil
}
