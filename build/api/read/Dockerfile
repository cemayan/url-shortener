FROM golang:1.19.2-alpine3.16 AS builder


WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/api/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/read-service  cmd/api/read/main.go

FROM scratch
COPY --from=builder /app/config/api/config-docker.yaml /app/config/api/config-docker.yaml
COPY --from=builder /app/read-service /app/read-service
EXPOSE 8081
ENTRYPOINT ["/app/read-service"]
