# Start from golang base image
FROM golang:1.23-alpine as builder

# Enable go modules
ENV GO111MODULE=on

# Install git and build tools
RUN apk update && apk add --no-cache git bash build-base

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ .

# Update go.mod if needed and build the application
RUN go mod tidy && CGO_ENABLED=1 go build -o ./main .

# Run executable
CMD ["./main"]