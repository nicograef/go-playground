FROM golang:1.24-alpine

RUN apk update

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /go/bin/queue ./main.go




