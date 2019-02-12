package main

import (
	"fmt"
	"os"
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

var (
	c    *aws.Config
	once sync.Once
)

func config() *aws.Config {
	once.Do(func() {
		c = &aws.Config{
			Region:   aws.String(REGION),
			Endpoint: aws.String(ENDPOINT),
		}
	})
	return c
}

func DB() *dynamodb.DynamoDB {
	return dynamodb.New(session.Must(session.NewSession(config())))
}

func main() {
	var i int
	for {
		fmt.Print("(1) Create DP Table:\n(2) Delete DP Table:\n(0) Quit\n")
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			fmt.Println(err)
		}
		switch i {
		case 1:
			CreateTable()
		case 2:
			DeleteTable()
		case 0:
			os.Exit(0)
		}
	}

}

func CreateTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ip"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("gid"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ip"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("gid"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(_const.DBTABLENAME),
	}

	result, err := DB().CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
	fmt.Println("DONE")
}

func DeleteTable() {
	res, err := DB().DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(_const.DBTABLENAME)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res)
	fmt.Println("DONE")
}
