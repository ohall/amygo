#!/bin/bash

docker stop $(docker ps -q)
go build
docker build -t amygo .
docker run -p 8000:8000 amygo