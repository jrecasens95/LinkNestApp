FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server ./server
EXPOSE 4000
CMD ["./server"]
