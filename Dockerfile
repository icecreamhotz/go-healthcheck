FROM golang:1.17-alpine

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download
RUN go get github.com/cosmtrek/air

COPY . .