FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/main.go


FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]
