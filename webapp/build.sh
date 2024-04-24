#!/bin/bash
docker build -t genai-go:demonext . 

docker images | grep genai-go
