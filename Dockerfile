FROM golang:1.25-alpine

WORKDIR /bot

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot ./cmd/bot/main.go

CMD ["./bot"]