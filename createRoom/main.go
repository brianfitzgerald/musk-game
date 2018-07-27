package main

import (
	"encoding/json"
	"math/rand"
	"musk-game/model"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoClient = dynamodb.New(sess)
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	disasters := []string{"World Hunger", "Obesity"}

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	chosenDisaster := disasters[rand.Intn(len(disasters))]

	playerNum, err := strconv.Atoi(request.QueryStringParameters["players"])

	if err != nil {
		panic(err)
	}

	room := model.Room{
		Disaster:  chosenDisaster,
		Code:      strconv.Itoa(rand.Intn(1000)),
		Players:   playerNum,
		MuskDrawn: false,
	}

	st, err := dynamodbattribute.MarshalMap(room)

	input := &dynamodb.PutItemInput{
		Item:      st,
		TableName: aws.String(model.TableName),
	}

	_, err = dynamoClient.PutItem(input)

	if err != nil {
		panic(err)
	}

	roomString, err := json.Marshal(room)

	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{Body: string(roomString), StatusCode: 200, Headers: map[string]string{
		"Access-Control-Allow-Origin":      "*",    // Required for CORS support to work
		"Access-Control-Allow-Credentials": "true", // Required for cookies, authorization headers with HTTPS
	}}, nil
}

func main() {
	lambda.Start(Handler)
}
