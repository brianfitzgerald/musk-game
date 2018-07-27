package main

import (
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Response struct {
	RoomCode string `json:"RoomCode"`
}

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

func Handler(request events.APIGatewayProxyRequest) (Response, error) {

	disasters := []string{"World Hunger", "Obesity"}

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	chosenDisaster := disasters[rand.Intn(len(disasters))]

	room := Room{
		Disaster:  chosenDisaster,
		Code:      string(rand.Intn(1000)),
		Players:   request.QueryStringParameters["players"],
		MuskDrawn: false,
	}

	st, err := dynamodbattribute.MarshalMap(room)

	input := &dynamodb.PutItemInput{
		Item:      st,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)

	if err != nil {
		panic(err)
	}

	return Response{
		RoomCode: room.Code,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
