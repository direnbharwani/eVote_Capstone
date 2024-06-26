package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"

	"github.com/direnbharwani/evote-capstone/app/server/common"
	chaincode "github.com/direnbharwani/evote-capstone/chaincode/src"
)

// ======================================================================================
// Lambda Definition
// ======================================================================================

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody LambdaRequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
		return errorResponse, nil
	}

	// Load default SDK configuration using Lambda's IAM role
	configuration, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	voterCredentialsTable := common.DynamoDBTable{
		TableName:    "voter-credentials",
		PartitionKey: "nric",
		SortKey:      "electionID",
	}
	if err = voterCredentialsTable.Init(configuration, true); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	voterCredentials, err := common.GetItem[common.VoterCredentials](ctx, &voterCredentialsTable, common.DynamoDBKeys{
		PartitonKeyValue: requestBody.NRIC,
		SortKeyValue:     requestBody.ElectionID,
	})
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}
	itemExists := (voterCredentials.NRIC != "" && voterCredentials.ElectionID != "")

	if itemExists { // Check if valid
		if voterCredentials.VoterID == "" || voterCredentials.BallotID == "" {
			errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s-%s has an invalid entry!", voterCredentials.NRIC, voterCredentials.ElectionID))
			return errorResponse, nil
		} else {
			errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, "Entry already exists!")
			return errorResponse, nil
		}
	}

	// Implicit: Item does not exist
	// Create new entry, generate voterID, register & enroll identity & create new ballot
	newVoterUUID, err := uuid.NewV7()
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to generate voter UUID: %v", err))
		return errorResponse, nil
	}
	voterID := "v-" + newVoterUUID.String()
	fmt.Printf("voterID: %s\n", voterID)

	// Register and enroll voter identity
	secret, err := registerIdentity(voterID)
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}
	if err = enrollIdentity(voterID, secret); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	// Create ballot
	newBallotUUID, err := uuid.NewV7()
	if err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to generate voter UUID: %v", err))
		return errorResponse, nil
	}
	ballotID := "b-" + newBallotUUID.String()
	fmt.Printf("ballotID: %s\n", ballotID)

	newBallot := chaincode.Ballot{
		Asset:      chaincode.Asset{ID: ballotID},
		ElectionID: requestBody.ElectionID,
		VoterID:    voterID,
	}

	if err = common.ChaincodeCreate(voterID, os.Getenv("KALEIDO_AUTH_TOKEN"), newBallot); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	// Put new item in dynamoDB
	newCredentials := common.VoterCredentials{
		NRIC:       requestBody.NRIC,
		ElectionID: requestBody.ElectionID,
		VoterID:    voterID,
		BallotID:   ballotID,
	}

	if err = common.PutItem[common.VoterCredentials](ctx, &voterCredentialsTable, newCredentials); err != nil {
		errorResponse := common.GenerateErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v", err))
		return errorResponse, nil
	}

	return common.GenerateSuccessResponse(""), nil
}

func main() {
	lambda.Start(Handler)
}

// =============================================================================
// Helpers
// =============================================================================

func registerIdentity(name string) (string, error) {
	endpoint := fmt.Sprintf("%s/identities", os.Getenv("KALEIDO_REST_API_ENDPOINT"))

	// Register

	requestBody, err := json.Marshal(map[string]string{
		"name": name,
		"type": "client",
	})
	if err != nil {
		return "", fmt.Errorf("failed to prepare register identity request: %v", err)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating register request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", os.Getenv("KALEIDO_AUTH_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending register request: %v", err)
	}
	defer response.Body.Close()

	responseBodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading register response body: %v", err)
	}
	// fmt.Println((string(responseBodyData)))

	var responseBody map[string]interface{}
	if err = json.Unmarshal(responseBodyData, &responseBody); err != nil {
		return "", fmt.Errorf("failed to parse register response body: %v", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("%v", responseBody["error"])
	}

	return responseBody["secret"].(string), nil
}

func enrollIdentity(name, secret string) error {
	endpoint := fmt.Sprintf("%s/identities/%s/enroll", os.Getenv("KALEIDO_REST_API_ENDPOINT"), name)

	requestBody, err := json.Marshal(map[string]string{
		"secret": secret,
	})
	if err != nil {
		return fmt.Errorf("failed to prepare enroll identity request: %v", err)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating enroll request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", os.Getenv("KALEIDO_AUTH_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending enroll request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		responseBodyData, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error reading enroll response body: %v", err)
		}

		var responseBody map[string]interface{}
		if err = json.Unmarshal(responseBodyData, &responseBody); err != nil {
			return fmt.Errorf("failed to parse register response body: %v", err)
		}
		return fmt.Errorf("%v", responseBody["error"])
	}

	return nil
}

// =============================================================================
// API Types
// =============================================================================

type LambdaRequestBody struct {
	NRIC       string `json:"NRIC"`
	ElectionID string `json:"ElectionID"`
}
