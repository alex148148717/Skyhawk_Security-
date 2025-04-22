FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN  go build -v -x -o app ./main.go


FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]