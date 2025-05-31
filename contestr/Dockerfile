FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd/app

FROM alpine:latest

WORKDIR /app

ENV CONFIG_PATH=/app/configs/config.yaml

COPY --from=builder /app/server .
COPY --from=builder /app/configs/config.yaml $CONFIG_PATH
RUN chmod +x server

EXPOSE 8080

ENTRYPOINT ["./server"]
