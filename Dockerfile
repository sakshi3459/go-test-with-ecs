## We specify the base image we need for our
## go application
FROM golang:1.12.0-alpine3.9

RUN apk add --no-cache python3 py3-pip \
    && pip3 install --upgrade pip \
    && pip3 install awscli \
    && rm -rf /var/cache/apk/*
    
## We create an /app directory within our
## image that will hold our application source
## files

RUN mkdir /app

## We copy everything in the root directory
## into our /app directory
ADD . /app

## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

## we run go build to compile the binary
## executable of our Go program
RUN apk add git
ADD go.mod go.sum ./
RUN go mod download
RUN go build -o main .
EXPOSE 6379
## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]
