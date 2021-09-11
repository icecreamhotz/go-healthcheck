FROM golang:1.16-alpine

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download
RUN go get github.com/cosmtrek/air

COPY . .