FROM golang AS builder

WORKDIR /src

COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main .

ENTRYPOINT  ["./main"]