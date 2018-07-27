package main

import (
	"fmt"
	"math/rand"
	"musk-game/model"

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

	roomCode := request.QueryStringParameters["code"]

	println("model.TableName", model.TableName)
	println("roomCode", roomCode)

	result, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(model.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Code": {
				S: aws.String(roomCode),
			},
		},
	})

	fmt.Println("oof result", result)
	fmt.Println("oof err", err)

	if err != nil {
		panic(err)
	}

	room := &model.Room{}

	fmt.Println("oof room", result)

	err = dynamodbattribute.UnmarshalMap(result.Item, room)

	if err != nil {
		panic(err)
	}

	muskDrawn := false

	if !room.MuskDrawn {
		fmt.Println("oof players", room.Players)
		playersAmt := int(room.Players)
		fmt.Println("oof playersAmt", playersAmt)
		chance := rand.Intn(playersAmt)
		fmt.Println("oof randomChance", chance)
		if chance == 0 {
			muskDrawn = true
			room.MuskDrawn = true
		}
	}

	room.Players = room.Players - 1

	updatedRoom, err := dynamodbattribute.MarshalMap(room)

	if err != nil {
		panic(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      updatedRoom,
		TableName: aws.String(model.TableName),
	}

	_, err = dynamoClient.PutItem(input)

	if err != nil {
		panic(err)
	}

	draw := room.Disaster

	if muskDrawn {
		draw = "Elon Musk"
	}

	return events.APIGatewayProxyResponse{Body: draw, StatusCode: 200, Headers: map[string]string{
		"Access-Control-Allow-Origin":      "*",    // Required for CORS support to work
		"Access-Control-Allow-Credentials": "true", // Required for cookies, authorization headers with HTTPS
	}}, nil
}

func main() {
	lambda.Start(Handler)
}
