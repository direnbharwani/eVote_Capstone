package common

import (
	"encoding/json"
	"fmt"

	paillier "github.com/direnbharwani/evote-capstone/paillier"
)

// Prettifies a JSON for output.
// Returns the prettified JSON as a plain string
func PrettyJSON(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func DecodeKeys(publicKeyBase64, privateKeyBase64 string) (*paillier.PublicKey, *paillier.PrivateKey, error) {
	publicKey, err := paillier.Base64Decode[paillier.PublicKey](publicKeyBase64)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding public key: %v", err)
	}

	privateKey, err := paillier.Base64Decode[paillier.PrivateKey](privateKeyBase64)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding private key: %v", err)
	}

	return publicKey, privateKey, nil
}
