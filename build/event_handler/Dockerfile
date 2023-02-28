FROM golang:1.19.2-alpine3.16 AS builder


WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/event_handler/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/event_handler  cmd/event_handler/main.go

FROM scratch
COPY --from=builder /app/config/event_handler/config-docker.yaml /app/config/event_handler/config-docker.yaml
COPY --from=builder /app/event_handler /app/event_handler
ENTRYPOINT ["/app/event_handler"]
