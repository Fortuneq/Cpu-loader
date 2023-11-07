FROM golang:latest AS builder


WORKDIR /app

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN  CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags="-s -w" -o main main.go


CMD ["./main"]



