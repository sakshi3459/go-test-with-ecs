## Specify the base image we need for our go application
FROM golang:1.17.0-alpine3.14

## go main application path
ENV SOURCE_APP /platform/service/daily-s3/*.go

## Create an /app directory within our
## image that will hold application source files
RUN mkdir /app


## Specify working directory for docker container
WORKDIR /app

## Copy application into /app directory
ADD $SOURCE_APP ./

## Download go packages and run go build to compile the binary
## executable of our Go program
RUN apk add git
ADD go.mod go.sum ./
## RUN go mod download
RUN go mod download github.com/aws/aws-sdk-go
RUN go mod download github.com/go-redis/redis/v7
RUN go mod download github.com/panjf2000/ants/v2
RUN go build -o main .

## Start command which kicks off created binary executable
CMD ["/app/main"]
