FROM golang:1.19.2-alpine3.16 AS builder


WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/api/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/write-service  cmd/api/write/main.go

FROM scratch
COPY --from=builder /app/config/api/config-docker.yaml /app/config/api/config-docker.yaml
COPY --from=builder /app/write-service /app/write-service
EXPOSE 8082
ENTRYPOINT ["/app/write-service"]
