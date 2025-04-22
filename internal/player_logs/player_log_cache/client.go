package player_log_cache

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Client struct {
	dynamodb  dynamodbiface.DynamoDBAPI
	tableName string
}

func New(tableName string, dynamodb dynamodbiface.DynamoDBAPI) *Client {
	c := &Client{tableName: tableName, dynamodb: dynamodb}
	return c
}

func (c *Client) PutItem(ctx context.Context, key string, value []byte) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
			"value": {
				B: value,
			},
		},
	}
	_, err := c.dynamodb.PutItem(input)
	return err
}
func (c *Client) GetItem(ctx context.Context, key string) ([]byte, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	}

	result, err := c.dynamodb.GetItem(input)
	if err != nil {
		return nil, err
	}

	valAttr, ok := result.Item["value"]
	if !ok || valAttr.B == nil {
		return nil, fmt.Errorf("key not found")
	}

	return valAttr.B, nil
}
