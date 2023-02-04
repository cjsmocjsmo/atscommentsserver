#!/bin/bash

sudo apt-get -y update && \
sudo apt-get -y dist-upgrade && \
sudo apt-get -y golang && \
sudo apt-get -y autoclean && \
sudo apt-get -y autoremove && \
mkdir ./log/logfile.txt && \ 
mkdir ./static && / #for pics
mkdir ./data && /
mkdir ./data/accepted && \ 
mkdir ./data/rejected && \
mkdir ./data/jailed


go get -v /go/src/atsflutter
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main /go/src/atsflutter