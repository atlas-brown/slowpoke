FROM golang:1.23-alpine

WORKDIR /app
COPY app/go.mod .
COPY app/go.sum .
COPY app/main.go .
RUN go build
