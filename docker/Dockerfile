FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 keeps the binary statically linked; -s -w drops the debug
# symbol table and DWARF data, which cuts roughly 30% off the image.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/main.go


FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]
