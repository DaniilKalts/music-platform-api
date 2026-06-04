FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o music-api ./cmd/api/main.go

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/music-api .
COPY api/v1 ./api/v1
COPY web/swagger ./web/swagger
CMD ["./music-api"]
