#!/bin/bash

# installing docker
sudo apt-get update
sudo apt-get install -y docker.io

# pulling the image
sudo docker pull mannuk24/challenge:latest

# stop and remove the existing container if it exists
sudo docker stop challenge
sudo docker rm challenge

sudo docker pull mannuk24/challenge:latest# running container
#
sudo docker run -d --name challenge -p 80:8050 mannuk24/challenge:latest
