FROM golang:1.23 AS builder

WORKDIR /PaymentSystem

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o payment-system ./cmd/main.go

CMD ["./payment-system"]