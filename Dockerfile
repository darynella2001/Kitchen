# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

LABEL maintainer="Andronovici Darinela"
WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./main/main.go

EXPOSE 8081

CMD ["./main"]
