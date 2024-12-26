FROM golang:1.23 AS builder

WORKDIR /NotificationSvc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification-svc ./cmd/main.go

CMD ["./notification-svc"]
