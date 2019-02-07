package db

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pokt-network/pocket-core/const"
)

// NOTE: You need to setup your AWS configuration
// first https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html

const (
	REGION   = "us-west-2"
	ENDPOINT = "http://localhost:8000"
)

type DB struct {
	dynamo *dynamodb.DynamoDB
	sync.Once
}

// "NewDB" returns a new database instance.
func NewDB() *DB {
	db := &DB{}
	db.init()
	return db
}

// "init" initializes the database.
func (db *DB) init() {
	var config *aws.Config
	db.Do(func() {
		config = &aws.Config{
			Region:   aws.String(_const.DBREIGON),
			Endpoint: aws.String(_const.DBENDPOINT),
		}
	})
	// start the session
	db.dynamo = dynamodb.New(session.Must(session.NewSession(config)))
}
