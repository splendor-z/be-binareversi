FROM golang:1.23-alpine

RUN apk add --no-cache git curl && \
    go install github.com/cosmtrek/air@v1.40.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]
