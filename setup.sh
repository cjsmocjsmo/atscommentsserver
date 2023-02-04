#!/bin/sh

sudo apt-get -y update && \
sudo apt-get -y dist-upgrade && \
sudo apt-get -y golang && \
sudo apt-get -y autoclean && \
sudo apt-get -y autoremove && \
mkdir ./log && \ 
mkdir ./static && / #for pics
mkdir ./data && /
mkdir ./data/accepted && \ 
mkdir ./data/admin && \
mkdir ./data/admin/profiles && \ 
mkdir ./data/backups && \ 
mkdir ./data/estcompleted && \ 
mkdir ./data/estimates && \ 
mkdir ./data/jailed && \ 
mkdir ./data/rejected && \

touch ./log/logfile.txt && \ 
touch ./data/admin/loggedinList.json && \ 

sudo mv ./ATSCommentsServer.service /etc/systemd/system/

go get -v /home/porthose_cjsmo_cjsmo/atscommentsserver/
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o atscommentsserver /home/porthose_cjsmo_cjsmo/atscommentsserver/

sudo systemctl start ATSCommentsServer