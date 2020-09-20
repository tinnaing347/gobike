FROM golang:1.15
RUN apt-get update && apt-get -y install --no-install-recommends apt-utils
RUN go get -v github.com/go-co-op/gocron

RUN mkdir /code 
WORKDIR /code/ 
ADD . /code/ 

RUN ls
RUN printenv 
