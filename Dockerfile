FROM golang AS builder

WORKDIR /src

COPY go.mod go.sum /
RUN go mod download

COPY . .

ENTRYPOINT  ["./main"]