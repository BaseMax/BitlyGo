FROM golang:1.19.0-alpine3.16

RUN apk update && apk add --no-cache bash

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]

RUN go mod download
RUN go install github.com/jackc/tern@latest

COPY . ./

RUN go build -o . ./...

EXPOSE 8000

CMD [ "./bitlygo" ]
