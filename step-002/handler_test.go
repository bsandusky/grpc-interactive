package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestValidateRequest(t *testing.T) {

	tests := []struct {
		in  registrationRequest
		err error
	}{
		{registrationRequest{"lou@gmail.com", "passphrase"}, nil},
		{registrationRequest{"", "passphrase"}, errors.New("Payload missing data: email_address required")},
		{registrationRequest{"lou@gmail.com", ""}, errors.New("Payload missing data: passphrase required")},
		{registrationRequest{"lou@gmail", "passphrase"}, errors.New("Payload data invalid: email_address not valid")},
	}

	for _, tt := range tests {
		err := validateRequest(tt.in)

		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Incorrect error value: got %v, want %v", err, tt.err)
		}
	}
}

func BenchmarkValidateRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validateRequest(registrationRequest{"lou@gmail.com", "passphrase"})
	}
}
