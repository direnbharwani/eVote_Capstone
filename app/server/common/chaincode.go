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

type InvokeType int

const (
	Transaction = iota + 1
	Query
)

// =============================================================================
// API Types
// =============================================================================

type ChaincodeInvocationHeaders struct {
	Type      string `json:"type,omitempty"`
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

func ChaincodeCreate[T chaincode.ITYPES](signer, authToken string, data T) error {
	function := fmt.Sprintf("Create%s", reflect.TypeOf(data).Name())

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = invokeChaincode(Transaction, signer, authToken, function, []string{string(rawData)}); err != nil {
		return err
	}

	return nil
}

// Queries a single object from the blockchain's world state
// Chaincode name, channel, and init are hardcoded
func ChaincodeQuery[T chaincode.ITYPES](signer, authToken, key string) (T, error) {
	var emptyObject T
	var result T

	function := fmt.Sprintf("Query%s", reflect.TypeOf(result).Name())

	chaincodeResponse, err := invokeChaincode(Query, signer, authToken, function, []string{key})
	if err != nil {
		return emptyObject, fmt.Errorf("%v", err)
	}

	// Temporary struct to convert the type accordingly
	type ChaincodeQueryRespondeBody struct {
		Headers map[string]interface{} `json:"headers"`
		Result  T                      `json:"result"`
	}

	var chaincodeResponseBody ChaincodeQueryRespondeBody
	err = json.Unmarshal(chaincodeResponse, &chaincodeResponseBody)
	if err != nil {
		return emptyObject, fmt.Errorf("error parsing chaincode response: %v", err)
	}

	return chaincodeResponseBody.Result, nil
}

// Queries a range of objects from the blockchain's world state
// Chaincode name, channel, and init are hardcoded
func ChaincodeQueryAll[T chaincode.ITYPES](signer, authToken string) ([]T, error) {
	var emptyObject T

	function := fmt.Sprintf("QueryAll%ss", reflect.TypeOf(emptyObject).Name())

	chaincodeResponse, err := invokeChaincode(Query, signer, authToken, function, []string{})
	if err != nil {
		return []T{}, fmt.Errorf("%v", err)
	}

	// Temporary struct to convert the type accordingly
	type ChaincodeQueryRespondeBody struct {
		Headers map[string]interface{} `json:"headers"`
		Result  []T                    `json:"result"`
	}

	var chaincodeResponseBody ChaincodeQueryRespondeBody
	err = json.Unmarshal(chaincodeResponse, &chaincodeResponseBody)
	if err != nil {
		return []T{}, fmt.Errorf("error parsing chaincode response: %v", err)
	}

	return chaincodeResponseBody.Result, nil
}

// =============================================================================
// Helpers
// =============================================================================

func invokeChaincode(invokeType InvokeType, signer, authToken, function string, args []string) ([]byte, error) {
	endpoint := os.Getenv("KALEIDO_REST_API_ENDPOINT")

	// Build chaincode request
	chaincodeInvocationHeaders := ChaincodeInvocationHeaders{
		Signer:    signer,
		Channel:   "default-channel",
		Chaincode: "evote_poc",
	}

	switch invokeType {
	case Transaction:
		{
			chaincodeInvocationHeaders.Type = "SendTransaction"
			endpoint += "/transactions"
		}
	case Query:
		{
			endpoint += "/query"
		}
	}

	chaincodeRequestBody := ChaincodeRequestBody{
		Headers: chaincodeInvocationHeaders,
		Func:    function,
		Args:    args,
		Init:    false,
	}

	chaincodeRequestJSONData, err := json.Marshal(chaincodeRequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare chaincode request body: %v", err)
	}

	fmt.Println(endpoint)
	fmt.Println(string(chaincodeRequestJSONData))

	chaincodeRequest, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(chaincodeRequestJSONData))
	if err != nil {
		return nil, fmt.Errorf("error creating chaincode request: %v", err)
	}

	chaincodeRequest.Header.Set("Content-Type", "application/json")
	chaincodeRequest.Header.Set("Authorization", authToken)

	// Invoke chaincode through REST API Gateway to Query State
	client := &http.Client{}
	chaincodeResponse, err := client.Do(chaincodeRequest)
	if err != nil {
		return nil, fmt.Errorf("error sending chaincode request: %v", err)
	}
	defer chaincodeResponse.Body.Close()

	// Parse response
	chaincodeResponseBodyData, err := io.ReadAll(chaincodeResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading chaincode response body: %v", err)
	}

	// Check for error in response
	if chaincodeResponse.StatusCode != 200 {
		var responseBody map[string]interface{}
		if err = json.Unmarshal(chaincodeResponseBodyData, &responseBody); err != nil {
			return nil, fmt.Errorf("%v", responseBody["error"])
		}
	}

	return chaincodeResponseBodyData, nil
}
