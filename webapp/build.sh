#!/bin/bash
docker build -t genai-go:demo . 

docker images | grep genai-go
