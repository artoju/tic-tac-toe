package db

import (
	"github.com/artoju/tic-tac-toe/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Init dynamoDB connection
func Init(conf *config.Config) (*dynamodb.DynamoDB, error) {
	creds := credentials.NewStaticCredentials(conf.DatabaseHandler.KeyId, conf.DatabaseHandler.SecretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      &conf.DatabaseHandler.Region,
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)
	return svc, nil
}
