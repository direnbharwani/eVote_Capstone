package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	chaincode "github.com/direnbharwani/evote-capstone/chaincode/src"
)

// =============================================================================
// API Types
// =============================================================================

type ChaincodeInvocationHeaders struct {
	Signer    string `json:"signer"`
	Channel   string `json:"channel"`
	Chaincode string `json:"chaincode"`
}

type ChaincodeRequestBody struct {
	Headers ChaincodeInvocationHeaders `json:"headers"`
	Func    string                     `json:"func"`
	Args    []string                   `json:"args"`
	Init    bool                       `json:"init"`
}

// =============================================================================
// Invocation Methods
// =============================================================================

// Queries a single object from the blockchain's world state
// Chaincode name, channel, and init are hardcoded
func ChaincodeQuery[T chaincode.ITYPES](signer, key string) (T, error) {
	var emptyObject T
	var result T

	// Build chaincode request
	chaincodeInvocationHeaders := ChaincodeInvocationHeaders{
		Signer:    signer,
		Channel:   "default-channel",
		Chaincode: "evote_poc",
	}

	chaincodeRequestBody := ChaincodeRequestBody{
		Headers: chaincodeInvocationHeaders,
		Func:    fmt.Sprintf("Query%s", reflect.TypeOf(result).Name()),
		Args:    []string{key},
		Init:    false,
	}

	chaincodeRequestJSONData, err := json.Marshal(chaincodeRequestBody)
	if err != nil {
		return emptyObject, fmt.Errorf("failed to prepare chaincode request body: %v", err)
	}

	chaincodeRequest, err := http.NewRequest("POST", "https://a0z8wc2w78-a0ve7t5vxf-connect.au0-aws-ws.kaleido.io/query", bytes.NewBuffer(chaincodeRequestJSONData))
	if err != nil {
		return emptyObject, fmt.Errorf("error creating chaincode request: %v", err)
	}

	chaincodeRequest.Header.Set("Content-Type", "application/json")
	chaincodeRequest.Header.Set("Authorization", os.Getenv("KALEIDO_AUTH_TOKEN"))

	// Invoke chaincode through REST API Gateway to Query State
	client := &http.Client{}
	chaincodeResponse, err := client.Do(chaincodeRequest)
	if err != nil {
		return emptyObject, fmt.Errorf("error sending chaincode request: %v", err)
	}
	defer chaincodeResponse.Body.Close()

	// Parse response
	chaincodeResponseBodyData, err := io.ReadAll(chaincodeResponse.Body)
	if err != nil {
		return emptyObject, fmt.Errorf("error chaincode reading response body: %v", err)
	}

	// Check for error in response
	if chaincodeResponse.StatusCode != 200 {
		var responseBody map[string]interface{}
		if err = json.Unmarshal(chaincodeResponseBodyData, &responseBody); err != nil {
			return emptyObject, fmt.Errorf("%v", responseBody["error"])
		}
	}

	var chaincodeResponseBody map[string]interface{}
	err = json.Unmarshal(chaincodeResponseBodyData, &chaincodeResponseBody)
	if err != nil {
		return emptyObject, fmt.Errorf("error parsing chaincode response: %v", err)
	}

	// Convert result field to T
	result = chaincodeResponseBody["result"].(T)
	return result, nil
}
