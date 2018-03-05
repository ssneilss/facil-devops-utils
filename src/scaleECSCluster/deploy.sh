#!/bin/bash
GOOS=linux go build -o main ./main.go
zip deployment.zip main
drone-lambda --region ap-northeast-1 \
  --function-name scaleECSCluster \
  --zip-file deployment.zip
rm deployment.zip && rm main
