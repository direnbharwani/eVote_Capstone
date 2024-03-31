package common

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBTYPES interface {
	VoterCredentials
}

type DynamoDBKeys struct {
	PartitonKeyValue string
	SortKeyValue     string
}

type DynamoDBTable struct {
	TableName       string
	PartitionKey    string
	SortKey         string
	sortKeyRequired bool
	client          *dynamodb.Client
}

func (t *DynamoDBTable) Init(config aws.Config, sortKeyRequired bool) error {
	t.sortKeyRequired = sortKeyRequired
	t.client = dynamodb.NewFromConfig(config)

	return nil
}

func (t *DynamoDBTable) CheckKeys() error {
	if t.TableName == "" {
		return fmt.Errorf("TableName is not set")
	}

	if t.PartitionKey == "" {
		return fmt.Errorf("PartitionKey is not set! At least the partitionKey must be set")
	}

	if t.sortKeyRequired && t.SortKey == "" {
		return fmt.Errorf("PartitionKey is not set! At least the partitionKey must be set")
	}

	return nil
}

// Gets an item from a dynamo DB Table. The keys must be specified
func GetItem[T DBTYPES](ctx context.Context, table *DynamoDBTable, keys DynamoDBKeys) (T, error) {
	var emptyObject T
	var result T

	if table.client == nil {
		return emptyObject, fmt.Errorf("%s has not been initialised", reflect.TypeOf(*table).Name())
	}

	if err := table.CheckKeys(); err != nil {
		return emptyObject, err
	}

	key := map[string]types.AttributeValue{
		table.PartitionKey: &types.AttributeValueMemberS{Value: keys.PartitonKeyValue},
	}
	if table.sortKeyRequired {
		key[table.SortKey] = &types.AttributeValueMemberS{Value: keys.SortKeyValue}
	}
	fmt.Println(key)

	getItemResult, err := table.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &table.TableName,
		Key:       key,
	})
	if err != nil {
		return emptyObject, fmt.Errorf("error getting item: %v", err)
	}

	if getItemResult.Item == nil {
		return emptyObject, nil
	}

	if err = attributevalue.UnmarshalMap(getItemResult.Item, &result); err != nil {
		return emptyObject, fmt.Errorf("error parsing result: %v", err)
	}

	return result, nil
}

func PutItem[T DBTYPES](ctx context.Context, table *DynamoDBTable, newItem T) error {
	if table.client == nil {
		return fmt.Errorf("%s has not been initialised", reflect.TypeOf(*table).Name())
	}

	newItemData, err := attributevalue.MarshalMap(newItem)
	if err != nil {
		return err
	}

	if _, err = table.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &table.TableName,
		Item:      newItemData,
	}); err != nil {
		return err
	}

	return nil
}

// =============================================================================
// DynamoDB Item Types
// =============================================================================

type VoterCredentials struct {
	NRIC       string `json:"NRIC" dynamodbav:"nric"`
	ElectionID string `json:"ElectionID" dynamodbav:"electionID"`
	VoterID    string `json:"VoterID" dynamodbav:"voterID"`
	BallotID   string `json:"BallotID" dynamodbav:"ballotID"`
}
