## We specify the base image we need for our
## go application
FROM golang:1.17.0-alpine3.14

# Add the keys
##ARG github_user
##ENV github_user=$github_user
##ARG github_personal_token
##ENV github_personal_token=$github_personal_token

## go main application path
##ENV SOURCE_APP /platform/service/daily-s3/*.*

## We create an /app directory within our
## image that will hold our application source
## files

RUN mkdir /app

## We copy everything in the root directory
## into our /app directory
#ADD . /app

## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

ADD . ./
ADD /platform/service/daily-s3/*.* ./

RUN ls
## we run go build to compile the binary
## executable of our Go program
RUN apk add git
## ADD go.mod go.sum ./
RUN cd /app/platform
RUN ls
## RUN cat /app/platform/service/daily-s3/main.go
##RUN cat /app/platform/service/daily-s3/go.mod

## RUN go mod download -dir platform/service/daily-s3/go.mod
## RUN go mod download github.com/aws/aws-sdk-go@v1.40.56
## RUN go mod download github.com/go-redis/redis/v7@v7.4.1
## RUN go mod download github.com/panjf2000/ants/v2@v2.4.6
##RUN go mod tidy
RUN go build -o main .
#EXPOSE 6379
## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]
