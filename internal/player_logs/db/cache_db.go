package db

import (
	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"skyhawk/internal/player_logs/config"
)

func NewCacheDB(config *config.Config) (dynamodbiface.DynamoDBAPI, error) {
	region := config.Dynamodb.Region
	endpoint := config.Dynamodb.Endpoint
	daxHostPorts := config.Dynamodb.DaxHostPorts
	isDax := config.Dynamodb.UseDax
	cfg := &aws.Config{
		Region: aws.String(region),
	}
	if endpoint != "" {
		cfg.Endpoint = aws.String(endpoint)
		cfg.Credentials = credentials.NewStaticCredentials("fake", "fake", "")
	}

	sess := session.Must(session.NewSession(cfg))
	var svc dynamodbiface.DynamoDBAPI
	var err error
	if isDax {
		daxCfg := dax.DefaultConfig()
		daxCfg.HostPorts = daxHostPorts
		daxCfg.Region = *cfg.Region

		svc, err = dax.New(daxCfg)
		if err != nil {
			return nil, err
		}

	} else {
		svc = dynamodb.New(sess)

	}

	return svc, nil

}
