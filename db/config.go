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

var (
	db     *Database
	dbOnce sync.Once
)

type Database struct {
	dynamo *dynamodb.DynamoDB
	sync.Mutex
}

// "DB" returns a new database instance.
func DB() *Database {
	dbOnce.Do(func() {
		db = &Database{}
		var config *aws.Config
		config = &aws.Config{
			Region:   aws.String(_const.DBREIGON),
			Endpoint: aws.String(_const.DBENDPOINT),
		}
		// start the session
		db.dynamo = dynamodb.New(session.Must(session.NewSession(config)))
	})
	return db
}
