FROM golang:1.23 AS builder

WORKDIR /HotelSvc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o hotel-service ./cmd/main.go

CMD ["./hotel-service"]
