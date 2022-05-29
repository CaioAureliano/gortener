#syntax=docker/dockerfile:1

FROM golang:1.17.8-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /gortener cmd/main.go

CMD ["/gortener"]