// This package is all persistent data storage related code.
package db

import (
	"fmt"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/node"
)

// "Add" 'puts' a node into the persistent data storage.
func (db *DB) Add(n node.Node) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(n)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(_const.Tablename),
	}
	res, err := db.dynamo.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return res, nil
}

// "Remove" 'deletes' a node from the persistent data storage.
func (db *DB) Remove(n node.Node) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"gid": {
				S: aws.String(n.GID),
			},
			"ip": {
				S: aws.String(n.IP),
			},
		},
		TableName: aws.String(_const.Tablename),
	}
	
	return db.dynamo.DeleteItem(input)
}

// "GetAll" returns all nodes from the database.
func (db *DB) GetAll() (*dynamodb.ScanOutput, error) {
	input := &dynamodb.ScanInput{TableName: aws.String(_const.Tablename)}
	return db.dynamo.Scan(input)
}
