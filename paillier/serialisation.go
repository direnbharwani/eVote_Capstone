package paillier

import (
	"encoding/base64"
	"encoding/json"
)

// =============================================================================
// JSON
// =============================================================================

// Serailises a key into bytes that can be represented as a JSON object
func SerializeToJSON[T ITYPES](key *T) ([]byte, error) {
	data, err := json.Marshal(*key)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// Deserialises a key from bytes that represented a JSON object
func DeserialiseFromJSON[T ITYPES](data []byte) (*T, error) {
	result := new(T)

	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}

// =============================================================================
// Base64
// =============================================================================

// Encodes a key as a base64 string
func Base64Encode[T ITYPES](key *T) (string, error) {
	data, err := json.Marshal(*key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// Decodes a key from a base64 string
func Base64Decode[T ITYPES](keyBase64 string) (*T, error) {
	result := new(T)

	data, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}
