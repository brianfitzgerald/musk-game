package main

import (
	"math/rand"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Room struct {
	Players   int
	MuskDrawn bool
	Code      string
	Disaster  string
}

const (
	tableName = "musk-game-rooms"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqsClient    = sqs.New(sess)
	dynamoClient = dynamodb.New(sess)
)

type Response struct {
	Draw string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {

	roomCode := request.QueryStringParameters["code"]

	result, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Code": {
				S: aws.String(roomCode),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	room := Room{}

	err = dynamodbattribute.UnmarshalMap(result.Item, room)

	if err != nil {
		panic(err)
	}

	muskDrawn := room.MuskDrawn

	if !room.MuskDrawn {
		chance := rand.Intn(room.Players)
		if chance == room.Players {
			muskDrawn = true
		}
	}

	room.MuskDrawn = muskDrawn
	room.Players = room.Players - 1

	updatedRoom, err := dynamodbattribute.MarshalMap(room)

	if err != nil {
		panic(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      updatedRoom,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)

	if err != nil {
		panic(err)
	}

	draw := room.Disaster

	if muskDrawn {
		draw = "Elon Musk"
	}

	return Response{
		Draw: draw,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
