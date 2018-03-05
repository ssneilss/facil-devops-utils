package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestDeployECS(t *testing.T) {
	t.Run("Should failed parsing json", func(t *testing.T) {
		_, err := parsePacakgeJSON("./none.json")
		if err == nil {
			t.Fail()
		}
	})

	t.Run("Should parse json version correctly", func(t *testing.T) {
		info, _ := parsePacakgeJSON("./package.json")
		if info.Name != "facil-admin" {
			t.Fail()
		}
		if info.Version != "0.1.0" {
			t.Fail()
		}
	})

	t.Run("Should get repository uri", func(t *testing.T) {
		URI, _ := getRepoURI("facil")
		if URI != "108327956849.dkr.ecr.ap-northeast-1.amazonaws.com/facil" {
			t.Fail()
		}
	})

	t.Run("Should parse task definition correctly", func(t *testing.T) {
		taskdef, _ := parseTaskDefinition("./taskdef.json", &taskVars{
			RepositoryURI: "REPO",
			Version:       "0.1.0",
		})
		if aws.StringValue(taskdef.ContainerDefinitions[0].Image) != "REPO:v_0.1.0" {
			t.Fail()
		}
	})
}
