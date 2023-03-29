FROM golang:1.20-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY . ./

RUN go mod download

EXPOSE 8080

RUN go build -o project

CMD ["./project"]

