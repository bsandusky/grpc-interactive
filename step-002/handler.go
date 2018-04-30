package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

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

func getDB() (*gorm.DB, error) {

	// Init the database
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=grpc-interactive password=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	return db, nil
}

func getUserToken(db *gorm.DB, req registrationRequest) (string, error) {

	// TODO: implement this!
	// check if user exists; return JWT if so
	// generate jwt and new user object
	// add user to db

	return "token", nil
}

func validateRequest(req registrationRequest) error {
	// check that payload is complete
	if req.EmailAddress == "" {
		return errors.New("Payload missing data: email_address required")
	} else if req.Passphrase == "" {
		return errors.New("Payload missing data: passphrase required")
	}

	// check that payload email_address is valid
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(req.EmailAddress) {
		return errors.New("Payload data invalid: email_address not valid")
	}
	// valid registrationRequest
	return nil
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// unmarshal request body
	var r registrationRequest
	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		log.Fatal(err)
	}

	if err := validateRequest(r); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	db, err := getDB()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	token, err := getUserToken(db, r)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// return jwt and instructions for usage + next step
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       token,
	}, nil
}

func main() {
	lambda.Start(handler)
}
