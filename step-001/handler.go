package main

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func getS3Data(svc s3iface.S3API, bucket, key string) (string, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	res, err := svc.GetObject(input)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf.String(), nil
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	svc := s3.New(session.New())
	data, _ := getS3Data(svc, "grpc-interactive", "step-001/data.json")
	headers := make(map[string]string)
	headers["content-type"] = "application/json"
	headers["step-two-url"] = os.Getenv("STEP_TWO_URL")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       data,
	}, nil
}

func main() {
	lambda.Start(handler)
}
