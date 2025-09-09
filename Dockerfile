FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /api ./cmd/api/main.go


FROM alpine:latest

WORKDIR /app

COPY ./configs/config.yaml ./configs/config.yaml

COPY --from=builder /api .

CMD ["./api"]

