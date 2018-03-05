package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Request struct {
	StackName      string `json:"stackname"`
	TargetCapacity int64  `json:"target"`
	MaxCapacity    int64  `json:"max"`
	MinCapacity    int64  `json:"min"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

var awsSession = session.New(
	&aws.Config{
		Region: aws.String("ap-northeast-1"),
	},
)

func Handler(request Request) (Response, error) {
	cloudformationClient := cloudformation.New(awsSession)
	result, cloudformationError := cloudformationClient.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(request.StackName),
	})
	if cloudformationError != nil {
		fmt.Printf("Failed to describe stacks on cloudformation stack %s", request.StackName)
		fmt.Println(cloudformationError)
		return Response{
			Message: fmt.Sprint("Failed"),
			Ok:      false,
		}, cloudformationError
	}
	outputs := result.Stacks[0].Outputs
	var spotFleetRequestID string
	for _, output := range outputs {
		if aws.StringValue(output.OutputKey) == "EcsSpotFleetRequestId" {
			spotFleetRequestID = aws.StringValue(output.OutputValue)
		}
	}

	autoScalingClient := applicationautoscaling.New(awsSession)
	_, autoScalingError := autoScalingClient.RegisterScalableTarget(&applicationautoscaling.RegisterScalableTargetInput{
		MaxCapacity:       aws.Int64(request.MaxCapacity),
		MinCapacity:       aws.Int64(request.MinCapacity),
		ResourceId:        aws.String(fmt.Sprintf("spot-fleet-request/%s", spotFleetRequestID)),
		ScalableDimension: aws.String("ec2:spot-fleet-request:TargetCapacity"),
		ServiceNamespace:  aws.String("ec2"),
	})
	if autoScalingError != nil {
		fmt.Printf("Failed to change autoScaling policy on spotFleetRequest %s:", spotFleetRequestID)
		fmt.Println(autoScalingError)
		return Response{
			Message: fmt.Sprint("Failed"),
			Ok:      false,
		}, autoScalingError
	}

	ec2Client := ec2.New(awsSession)
	_, spotFleetRequestError := ec2Client.ModifySpotFleetRequest(&ec2.ModifySpotFleetRequestInput{
		SpotFleetRequestId: aws.String(spotFleetRequestID),
		TargetCapacity:     aws.Int64(request.TargetCapacity),
	})
	if spotFleetRequestError != nil {
		fmt.Printf("Failed to change spotFleetRequest %s to capacity %d:", spotFleetRequestID, request.TargetCapacity)
		fmt.Println(spotFleetRequestError)
		return Response{
			Message: fmt.Sprint("Failed"),
			Ok:      false,
		}, spotFleetRequestError
	}

	return Response{
		Message: fmt.Sprintf("Successfully change spotFleetRequest on %s to %d", request.StackName, request.TargetCapacity),
		Ok:      true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
