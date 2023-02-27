FROM golang:1.20-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./
COPY docker/init.sql /docker-entrypoint-initdb.d/init.sql

ENV APP_PORT=8080

EXPOSE 8080

RUN go build -o project

CMD ["./project"]

