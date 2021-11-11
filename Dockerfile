# syntax=docker/dockerfile:1

FROM golang:1.17.2-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /goapp-rest-api

EXPOSE 8080

CMD [ "/goapp-rest-api" ]