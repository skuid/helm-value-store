package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/skuid/helm-value-store/store"
)

// ReleaseStore stores and retrieves releases from a DynamoDB table
type ReleaseStore struct {
	tableName string
	sess      *session.Session
}

// NewReleaseStore Creates a new ReleaseStore
func NewReleaseStore(tableName string) (store.ReleaseStore, error) {
	rs := &ReleaseStore{tableName: tableName}

	sess, err := session.NewSession(
		&aws.Config{CredentialsChainVerboseErrors: aws.Bool(true)},
	)
	if err != nil {
		return nil, err
	}
	rs.sess = sess

	return rs, nil
}

// Get gets a release by it's UniqueID
func (rs ReleaseStore) Get(ctx context.Context, uniqueID string) (*store.Release, error) {
	svc := dynamodb.New(rs.sess)

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UniqueID": {
				S: aws.String(uniqueID),
			},
		},
		TableName:      aws.String(rs.tableName),
		ConsistentRead: aws.Bool(true),
	}
	resp, err := svc.GetItemWithContext(ctx, params)
	if err != nil {
		return nil, err
	}
	avm := attributeValueMap(resp.Item)
	return avm.MarshalRelease()

}

// Delete deletes a release by it's UniqueID
func (rs ReleaseStore) Delete(ctx context.Context, uniqueID string) error {
	svc := dynamodb.New(rs.sess)

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UniqueID": {
				S: aws.String(uniqueID),
			},
		},
		TableName: aws.String(rs.tableName),
	}
	_, err := svc.DeleteItemWithContext(ctx, params)
	return err
}

// Put creates or updates a release in DynamoDB
func (rs ReleaseStore) Put(ctx context.Context, r store.Release) error {
	svc := dynamodb.New(rs.sess)

	avm := attributeValueMap{}
	err := avm.UnmarshalRelease(r)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      avm,
		TableName: aws.String(rs.tableName),
	}
	_, err = svc.PutItemWithContext(ctx, params)
	return err
}

// List returns releases from DynamoDB
func (rs ReleaseStore) List(ctx context.Context, selector map[string]string) (store.Releases, error) {
	svc := dynamodb.New(rs.sess)

	// Dynamo doesn't support indexes on map types
	params := &dynamodb.ScanInput{
		TableName:      aws.String(rs.tableName),
		ConsistentRead: aws.Bool(true),
		Select:         aws.String("ALL_ATTRIBUTES"),
	}
	resp, err := svc.ScanWithContext(ctx, params)

	if err != nil {
		return nil, err
	}

	response := store.Releases{}
	for _, item := range resp.Items {
		avm := attributeValueMap(item)
		r, _ := avm.MarshalRelease()
		if r.MatchesSelector(selector) {
			response = append(response, *r)
		}
	}
	return response, nil
}

// Load bulk-writes releases to DynamoDB
func (rs ReleaseStore) Load(ctx context.Context, releases store.Releases) error {
	svc := dynamodb.New(rs.sess)
	params := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{},
	}
	writeRequests := []*dynamodb.WriteRequest{}

	for i, r := range releases {
		avm := attributeValueMap{}
		err := avm.UnmarshalRelease(r)
		if err != nil {
			return err
		}
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: avm},
		})

		params.RequestItems[rs.tableName] = writeRequests
		if i%25 == 0 {
			_, err := svc.BatchWriteItemWithContext(ctx, params)
			if err != nil {
				fmt.Println(err)
				return err
			}
			params = &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{},
			}
			writeRequests = []*dynamodb.WriteRequest{}
		}
	}
	if len(writeRequests) > 0 {
		_, err := svc.BatchWriteItemWithContext(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// Setup creates the table in DynamoDB if it doesn't exist. This call waits on
// the creation of the table to return
func (rs ReleaseStore) Setup(ctx context.Context) error {
	if !rs.tableExists(ctx) {
		err := rs.createTable(ctx)
		if err != nil {
			return err
		}
		svc := dynamodb.New(rs.sess)
		err = svc.WaitUntilTableExistsWithContext(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(rs.tableName)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs ReleaseStore) createTable(ctx context.Context) error {
	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("UniqueID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("UniqueID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(rs.tableName),
	}
	svc := dynamodb.New(rs.sess)
	_, err := svc.CreateTableWithContext(ctx, params)
	return err
}

func (rs ReleaseStore) tableExists(ctx context.Context) bool {
	svc := dynamodb.New(rs.sess)
	_, err := svc.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(rs.tableName),
	})
	if err != nil {
		return false
	}
	return true
}
