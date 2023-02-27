FROM golang:1.20-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY . ./

RUN go mod download

ENV APP_PORT=8080

EXPOSE 8080

RUN go build -o project

CMD ["./project"]

