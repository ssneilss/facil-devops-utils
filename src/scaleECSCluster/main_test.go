package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		request Request
		expect  string
		err     error
	}{
		{
			request: Request{
				StackName:      "EC2ContainerService-jenkins-slaves",
				TargetCapacity: 2,
				MaxCapacity:    10,
				MinCapacity:    2,
			},
			expect: "Successfully change spotFleetRequest on EC2ContainerService-jenkins-slaves to 2",
			err:    nil,
		},
		{
			request: Request{
				StackName:      "none",
				TargetCapacity: 2,
				MaxCapacity:    10,
				MinCapacity:    2,
			},
			expect: "Failed",
			err:    nil,
		},
	}

	for _, test := range tests {
		response, _ := Handler(test.request)
		assert.Equal(t, test.expect, response.Message)
	}
}
