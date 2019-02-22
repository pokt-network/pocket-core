package db

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pokt-network/pocket-core/config"
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
		con := config.GlobalConfig()
		db = &Database{}
		var c *aws.Config
		c = &aws.Config{
			Region:   aws.String(config.GlobalConfig().DBRegion),
			Endpoint: aws.String(con.DBEndpoint),
		}
		// start the session
		db.dynamo = dynamodb.New(session.Must(session.NewSession(c)))
	})
	return db
}
