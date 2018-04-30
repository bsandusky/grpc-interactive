package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type registrationRequest struct {
	EmailAddress string `json:"email_address"`
	Passphrase   string `json:"passphrase"`
}

// User object stores basic user information
type User struct {
	gorm.Model
	UUID         uuid.UUID `json:"uuid"`
	EmailAddress string    `json:"email_address"`
	Passphrase   string    `json:"passphrase"`
	Token        string    `json:"token"`
}

func registerUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// unmarshal request body
	var r registrationRequest
	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		log.Fatal(err)
	}

	// check that payload is complete
	if r.EmailAddress == "" || r.Passphrase == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       errors.New("Payload missing data").Error(),
		}, nil
	}

	// TODO: implement this!
	// check that payload is valid
	if false {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       errors.New("Payload data invalid").Error(),
		}, nil
	}

	// check if user exists; return JWT if so
	// generate jwt and new user object
	// add user to db

	// return jwt and instructions for usage + next step
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Hi",
	}, nil
}

func main() {
	lambda.Start(registerUser)
}
