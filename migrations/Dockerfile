FROM golang:1.21

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.18.0

WORKDIR /app
COPY . .

CMD ["goose", "-dir", ".", "postgres", "postgres://user:pass@postgres:5444/player_stats?sslmode=disable", "up"]
