package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
)

type dynamoDBRepo struct {
	tableName string
}

// NewDynamoDBRepository is the constructor function for the repo
func NewDynamoDBRepository() PostRepository {
	return &dynamoDBRepo{
		tableName: "posts",
	}
}

func createDynamoDBClient() *dynamodb.DynamoDB {
	// creates a new session by using .aws/credentials and .aws/config (gets credentials from local env)
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	return dynamodb.New(sess)
}

// Save ...
func (repo *dynamoDBRepo) Save(post *entity.Post) (*entity.Post, error) {
	// Get a DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	// Transform the post to map[string]*dynamodb.AttributeValue
	attributeValue, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return nil, err
	}

	// Create the Item Input
	item := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(repo.tableName),
	}

	// Save the Item into DynamoDB
	_, err = dynamoDBClient.PutItem(item)
	if err != nil {
		return nil, err
	}

	// Return the post
	return post, nil
}

func (repo *dynamoDBRepo) FindAll() ([]entity.Post, error) {

	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	// Make the DynamoDB Query API call
	result, err := dynamoDBClient.Scan(params)
	if err != nil {
		return nil, err
	}

	// Create a psots array and add all the existing posts
	var posts []entity.Post = []entity.Post{}
	for _, i := range result.Items {
		post := entity.Post{}
		// unmarshal item
		err = dynamodbattribute.UnmarshalMap(i, &post)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *dynamoDBRepo) FindOne(id string) (*entity.Post, error) {
	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	// Get the item by ID
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	// Map the dynamodb element to the post structure
	post := entity.Post{}

	// unmarshal item
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		panic(err)
	}

	// Return the pointer to the post
	return &post, nil
}

func (repo *dynamoDBRepo) Delete(post *entity.Post) error {

	return nil
}
