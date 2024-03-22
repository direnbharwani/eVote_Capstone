package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	paillier "github.com/direnbharwani/go-paillier/pkg"
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
	publicKeyData, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding public key: %v", err)
	}

	privateKeyData, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding private key: %v", err)
	}

	publicKey, err := paillier.DeserialiseJSON[paillier.PublicKey](publicKeyData)
	if err != nil {
		return nil, nil, fmt.Errorf("error unparsing public key: %v", err)
	}

	privateKey, err := paillier.DeserialiseJSON[paillier.PrivateKey](privateKeyData)
	if err != nil {
		return nil, nil, fmt.Errorf("error unparsing private key: %v", err)
	}

	return publicKey, privateKey, nil
}
