package chaincode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"time"

	paillier "github.com/direnbharwani/go-paillier/pkg"
)

// ITYPES is a union set type constraint
// that enforces only allowable types are passed to smart contract methods.
type ITYPES interface {
	Ballot | Candidate | Election

	Type() string
	Validate() error
	IsEqual(other interface{}) bool
}

type Asset struct {
	ID string `json:"ID"`
}

// =============================================================================
// Errors
// =============================================================================

type ObjectValidationError struct {
	ErrorMessage string
	ObjectType   string
}

func (e *ObjectValidationError) Error() string {
	return fmt.Sprintf("%s is invalid! %s", e.ObjectType, e.ErrorMessage)
}

type ObjectEqualityError struct {
	Key        string
	ObjectType string
}

func (e *ObjectEqualityError) Error() string {
	return fmt.Sprintf("%s %s have identical states!", e.ObjectType, e.Key)
}

type CompositeKeyCreationError struct {
	ErrorMessage string
	Key          string
	ObjectType   string
}

func (e *CompositeKeyCreationError) Error() string {
	return fmt.Sprintf("unable to create composite key for %s %s: %s", e.ObjectType, e.Key, e.ErrorMessage)
}

type WorldStateInteractionError struct {
	ErrorMessage string
	Key          string
}

func (e *WorldStateInteractionError) Error() string {
	return fmt.Sprintf("unable to interact with world state for %s: %s", e.Key, e.ErrorMessage)
}

type WorldStateReadFailureError struct {
	Key string
}

func (e *WorldStateReadFailureError) Error() string {
	return fmt.Sprintf("cannot read world state with key %s", e.Key)
}

// =============================================================================
// Election
// =============================================================================

// Defines an election
// Asset ID for Elections are prefixed with e-
type Election struct {
	Asset      Asset    `json:"Asset"`
	Candidates []string `json:"Candidates"`
	EndTime    string   `json:"EndTime"`
	Name       string   `json:"Name"`
	StartTime  string   `json:"StartTime"`
}

func (e Election) Type() string {
	return reflect.TypeOf(e).String()
}

func (e Election) Validate() error {
	objectType := reflect.TypeOf(e).String()

	if e.Asset.ID == "" {
		return &ObjectValidationError{"missing ID", objectType}
	}

	startTime, err := time.Parse(time.DateTime, e.StartTime)
	if err != nil {
		return &ObjectValidationError{err.Error(), objectType}
	}

	endTime, err := time.Parse(time.DateTime, e.EndTime)
	if err != nil {
		return &ObjectValidationError{err.Error(), objectType}
	}

	// EndTime must be after StartTime
	if endTime.Before(startTime) {
		return &ObjectValidationError{"EndTime must be after StartTime", objectType}
	}

	return nil
}

func (e Election) IsEqual(other interface{}) bool {
	otherObj, ok := other.(Election)
	if !ok {
		return false
	}

	if e.Asset != otherObj.Asset {
		return false
	}

	// Check if Candidates slices are equal
	if len(e.Candidates) != len(otherObj.Candidates) {
		return false
	}
	for i := range e.Candidates {
		if e.Candidates[i] != otherObj.Candidates[i] {
			return false
		}
	}

	// Check other fields for equality
	if e.EndTime != otherObj.EndTime || e.Name != otherObj.Name || e.StartTime != otherObj.StartTime {
		return false
	}

	return true
}

func (e Election) IsActive() bool {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().In(loc)

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

// Defines a electoral candidate with a public key for encrypting the count.
// The private key is omitted such that the count cannot be decrypted.
// Asset ID for Candidates are prefixed with c-
type Candidate struct {
	Asset      Asset    `json:"Asset"`
	Count      *big.Int `json:"Count"`
	ElectionID string   `json:"ElectionID"`
	Name       string   `json:"Name"`
	PublicKey  string   `json:"PublicKey"`
}

func (c Candidate) Type() string {
	return reflect.TypeOf(c).String()
}

func (c Candidate) Validate() error {
	objectType := reflect.TypeOf(c).String()

	if c.Asset.ID == "" {
		return &ObjectValidationError{"missing ID", objectType}
	}

	if c.ElectionID == "" {
		return &ObjectValidationError{"missing ElectionID", objectType}
	}

	return nil
}

func (c Candidate) IsEqual(other interface{}) bool {
	otherObj, ok := other.(Candidate)
	if !ok {
		return false
	}

	if c.Asset != otherObj.Asset {
		return false
	}

	if c.Count.Cmp(otherObj.Count) != 0 {
		return false
	}

	if c.PublicKey != otherObj.PublicKey || c.ElectionID != otherObj.ElectionID || c.Name != otherObj.Name {
		return false
	}

	return true
}

func (c *Candidate) Init() error {
	if c.PublicKey == "" {
		errorMessage := fmt.Sprintf("candidate %s is missing a public key! Unable to initialise", c.Asset.ID)
		return errors.New(errorMessage)
	}

	var err error

	decodedPublicKeyJSON, err := base64.StdEncoding.DecodeString(c.PublicKey)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to decode public key for candidate %s! Unable to initialise", c.Asset.ID)
		return errors.New(errorMessage)
	}

	publicKey, err := paillier.DeserialiseJSON[paillier.PublicKey](decodedPublicKeyJSON)
	if err != nil {
		return err
	}

	c.Count, err = paillier.Encrypt(publicKey, big.NewInt(0))
	if err != nil {
		return err
	}

	return nil
}

func (c *Candidate) IncrementCount() error {
	if c.PublicKey == "" {
		errorMessage := fmt.Sprintf("candidate %s is missing a public key! Unable to modify count", c.Asset.ID)
		return errors.New(errorMessage)
	}

	decodedPublicKeyJSON, err := base64.StdEncoding.DecodeString(c.PublicKey)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to decode public key for candidate %s! Unable to initialise", c.Asset.ID)
		return errors.New(errorMessage)
	}

	publicKey, err := paillier.DeserialiseJSON[paillier.PublicKey](decodedPublicKeyJSON)
	if err != nil {
		return err
	}

	c.Count = paillier.AddEncryptedWithPlain(publicKey, c.Count, big.NewInt(1))

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
func (b Ballot) Validate() error {
	objectType := reflect.TypeOf(b).String()

	if b.Asset.ID == "" {
		return &ObjectValidationError{"missing ID", objectType}
	}

	if b.ElectionID == "" {
		return &ObjectValidationError{"missing ElectionID", objectType}
	}

	return nil
}

func (b Ballot) IsEqual(other interface{}) bool {
	otherObj, ok := other.(Ballot)
	if !ok {
		return false
	}

	if b.Asset != otherObj.Asset {
		return false
	}

	// Check if Candidates slices are equal
	if len(b.Candidates) != len(otherObj.Candidates) {
		return false
	}
	for i := range b.Candidates {
		if b.Candidates[i].IsEqual(otherObj.Candidates[i]) {
			return false
		}
	}

	// Check other fields for equality
	if b.ElectionID != otherObj.ElectionID || b.VoterID != otherObj.VoterID || b.Voted != otherObj.Voted {
		return false
	}

	return true
}

func (b *Ballot) Vote(candidateID string) error {
	if b.Voted {
		errorMessage := fmt.Sprintf("ballot %s has already been cast! unable to vote", b.Asset.ID)
		return errors.New(errorMessage)
	}

	candidateFound := false
	for i, c := range b.Candidates {
		if c.Asset.ID == candidateID {
			candidateFound = true

			b.Candidates[i].IncrementCount()
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
