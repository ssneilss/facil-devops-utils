package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type pacakgeJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type taskVars struct {
	RepositoryURI string
	Version       string
}

var awsSession = session.New(
	&aws.Config{
		Region: aws.String("ap-northeast-1"),
	},
)

var client = ecs.New(awsSession)

func parsePacakgeJSON(path string) (pacakgeJSON, error) {
	var info pacakgeJSON
	content, err := ioutil.ReadFile(path)
	json.Unmarshal(content, &info)
	return info, err
}

func getRepoURI(name string) (string, error) {
	ecrClient := ecr.New(awsSession)
	output, err := ecrClient.DescribeRepositories(&ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(name)},
	})
	URI := aws.StringValue(output.Repositories[0].RepositoryUri)
	return URI, err
}

func deployImage(URI string, Version string) {
	repo := fmt.Sprintf("%s:v_%s", URI, Version)
	// Build docker image
	build := exec.Command("docker", "build", "-t", repo, ".")
	var buildOutput bytes.Buffer
	build.Stdout = &buildOutput
	if buileErr := build.Run(); buileErr != nil {
		log.Fatal(buileErr)
	}
	log.Println(buildOutput.String())
	// Get Login command
	getLogin := exec.Command("aws", "ecr", "get-login", "--region", "ap-northeast-1", "--no-include-email")
	var getLoginOuput bytes.Buffer
	getLogin.Stdout = &getLoginOuput
	if getLoginErr := getLogin.Run(); getLoginErr != nil {
		log.Fatal(getLoginErr)
	}
	log.Println(getLoginOuput.String())
	// Login to ECR
	loginCommand := strings.Fields(getLoginOuput.String())
	var loginOuput bytes.Buffer
	login := exec.Command(loginCommand[0], loginCommand[1:]...)
	login.Stdout = &loginOuput
	if loginErr := login.Run(); loginErr != nil {
		log.Fatal(loginErr)
	}
	log.Println(loginOuput.String())
	// Push image to ECR Repository
	push := exec.Command("docker", "push", repo)
	var pushOuput bytes.Buffer
	push.Stdout = &pushOuput
	if pushErr := push.Run(); pushErr != nil {
		log.Fatal(pushErr)
	}
	log.Println(pushOuput.String())
}

func parseTaskDefinition(path string, v *taskVars) (*ecs.TaskDefinition, error) {
	var taskdef *ecs.TaskDefinition
	content, err := ioutil.ReadFile(path)

	var str string
	str = string(content)
	str = strings.Replace(str, "%REPOSITORY_URI%", v.RepositoryURI, -1)
	str = strings.Replace(str, "%VERSION%", v.Version, -1)

	json.Unmarshal([]byte(str), &taskdef)

	return taskdef, err
}

func registerTaskDefinition(taskDefinition *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {
	output, err := client.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    taskDefinition.ContainerDefinitions,
		Cpu:                     taskDefinition.Cpu,
		ExecutionRoleArn:        taskDefinition.ExecutionRoleArn,
		Family:                  taskDefinition.Family,
		Memory:                  taskDefinition.Memory,
		NetworkMode:             taskDefinition.NetworkMode,
		PlacementConstraints:    taskDefinition.PlacementConstraints,
		RequiresCompatibilities: taskDefinition.RequiresCompatibilities,
		TaskRoleArn:             taskDefinition.TaskRoleArn,
		Volumes:                 taskDefinition.Volumes,
	})
	fmt.Println(output, err)
	return output.TaskDefinition, err
}

func updateService(input *ecs.UpdateServiceInput, taskdef *ecs.TaskDefinition) {
	result, _ := client.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  input.Cluster,
		Services: []*string{input.Service},
	})
	if len(result.Services) > 0 {
		log.Println("Start updating service")
		_, updateErr := client.UpdateService(input)
		if updateErr != nil {
			log.Println(updateErr)
		} else {
			log.Println("Updated service successfully")
		}
	} else {
		log.Println("Start creating service")
		_, err := client.CreateService(&ecs.CreateServiceInput{
			Cluster:        input.Cluster,
			ServiceName:    input.Service,
			TaskDefinition: input.TaskDefinition,
			DesiredCount:   aws.Int64(1),
		})
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Created service successfully")
		}
	}
}

func main() {
	var (
		Cluster         string
		Service         string
		PacakgeJSONPath string
		TaskdefPath     string
	)

	flag.StringVar(&Cluster, "cluster", "", "ECS Cluster Name")
	flag.StringVar(&Service, "service", "", "ECS Service Name")
	flag.StringVar(&PacakgeJSONPath, "package-json", "package.json", "Path to package.json")
	flag.StringVar(&TaskdefPath, "taskdef", "taskdef.json", "Path to taskdef.json")
	flag.Parse()

	info, _ := parsePacakgeJSON(PacakgeJSONPath)
	repoURI, _ := getRepoURI(info.Name)
	log.Println("Deployment Info:", info, repoURI)

	deployImage(repoURI, info.Version)
	log.Printf("Deployed image to %s_v:%s", repoURI, info.Version)

	taskdef, _ := parseTaskDefinition(TaskdefPath, &taskVars{
		RepositoryURI: repoURI,
		Version:       info.Version,
	})
	log.Println("Taskdef", taskdef)

	registeredTask, _ := registerTaskDefinition(taskdef)

	updateService(&ecs.UpdateServiceInput{
		Cluster:        aws.String(Cluster),
		Service:        aws.String(Service),
		TaskDefinition: registeredTask.TaskDefinitionArn,
	}, taskdef)
}
