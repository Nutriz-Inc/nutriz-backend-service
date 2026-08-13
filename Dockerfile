FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server main.go
RUN go build -o migrate cmd/migrations/run.go cmd/migrations/seed.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/.env.development .
COPY --from=builder /app/shared/database/migrations ./shared/database/migrations

EXPOSE 3333

CMD ["./server"]