FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
RUN chmod +x server

EXPOSE 8080

ENTRYPOINT ["./server"]
