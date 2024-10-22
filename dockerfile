FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download;
RUN go mod tidy;
RUN go build -o /app/main

