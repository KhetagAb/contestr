FROM golang:1.23-alpine AS builder

# RUN apk add --no-cache bash git curl make gcc musl-dev
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/google/wire/cmd/wire@v0.6.0

COPY . .

RUN sh ./scripts/codegen.sh

RUN CGO_ENABLED=0 \
    go build -a -installsuffix cgo -ldflags="-w -s" \
    -o server ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

ENV CONFIG_PATH=/app/configs/config.yaml

COPY --from=builder /app/server .
COPY --from=builder /app/configs/config.yaml $CONFIG_PATH
RUN chmod +x server

EXPOSE 8080

ENTRYPOINT ["./server"]
