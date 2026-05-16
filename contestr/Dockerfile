FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
RUN go install github.com/google/wire/cmd/wire@v0.6.0

COPY . .

RUN sh ./scripts/codegen.sh

RUN CGO_ENABLED=0 \
    go build -a -installsuffix cgo -ldflags="-w -s" \
    -o server ./cmd/app

FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

ENV CONFIG_PATH=/app/configs/config.yaml

COPY --from=builder /app/server .
COPY --from=builder /app/configs/config.yaml $CONFIG_PATH
RUN chmod +x server

EXPOSE 8080

ENTRYPOINT ["./server"]
