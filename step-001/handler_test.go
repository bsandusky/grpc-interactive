package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Mock S3 for testing
type mockS3 struct {
	s3iface.S3API
}

// Mock io.ReadCloser for testing s3.GetObjectOutput.Body
type nopCloser struct {
	io.Reader
}

// Nil implement Close() to satisfy io.ReadCloser interface
func (nopCloser) Close() error { return nil }

// Implement GetObject from S3iface library with testing logic
func (m mockS3) GetObject(in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {

	if !reflect.DeepEqual(in.Bucket, aws.String("grpc-interactive")) {
		return &s3.GetObjectOutput{}, awserr.New("InvalidBucketName", "The specified bucket is not valid.", nil)

	} else if !reflect.DeepEqual(in.Key, aws.String("step-001/data.json")) {
		return &s3.GetObjectOutput{}, awserr.New("NoSuchKey", "The resource you requested does not exist", nil)
	}

	return &s3.GetObjectOutput{Body: nopCloser{bytes.NewBufferString("some body data")}}, nil
}

func TestGetS3Data(t *testing.T) {
	mock := &mockS3{}

	t.Run("get *s3.GetObjectOutput when passing valid bucket/key combination", func(t *testing.T) {
		got, _ := getS3Data(mock, "grpc-interactive", "step-001/data.json")
		want := "some body data"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("incorrect output: got %s, want %s", got, want)
		}
	})

	t.Run("get error when passing invalid bucket", func(t *testing.T) {
		_, got := getS3Data(mock, "incorrect_bucket", "step-001/data.json")
		want := awserr.New("InvalidBucketName", "The specified bucket is not valid.", nil)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("incorrect error: got %s, want %s", got, want)
		}
	})

	t.Run("get error when passing invalid key", func(t *testing.T) {
		_, got := getS3Data(mock, "grpc-interactive", "invalid_key")
		want := awserr.New("NoSuchKey", "The resource you requested does not exist", nil)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("incorrect error: got %s, want %s", got, want)
		}
	})
}
