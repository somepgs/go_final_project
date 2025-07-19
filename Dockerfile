FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo ./cmd/main.go

ENV TODO_PORT=7540 \
    TODO_DBFILE=/data/scheduler.db \
    TODO_PASSWORD=12345

EXPOSE ${TODO_PORT}

CMD ["/todo"]