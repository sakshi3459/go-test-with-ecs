## Specify the base image we need for our go application
FROM golang:1.17.0-alpine3.14

## Create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app

## Copy everything in the root directory
## into our /app directory
ADD . /app

## Specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app

## Download go packages and run go build to compile the binary
## executable of our Go program
RUN apk add git
ADD /platform/service/daily-s3/*.go ./
RUN ls /app
ADD go.mod go.sum ./
RUN go mod download

RUN go build -o main .

## Start command which kicks off created binary executable
CMD ["/app/main"]
