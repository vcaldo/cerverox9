FROM golang:1.23.4 AS builder

WORKDIR /bot

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 go build -o /discord_bot  ./cmd

FROM alpine:3.21

COPY --from=builder /discord_bot /discord_bot

CMD ["/discord_bot"]