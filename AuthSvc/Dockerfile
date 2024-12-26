FROM golang:1.23 AS builder

WORKDIR /AuthSvc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

CMD ["./auth-service"]