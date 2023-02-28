FROM golang:1.19.2-alpine3.16 AS builder


WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/url_encoder/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/url_encoder  cmd/url_encoder/main.go

FROM scratch
COPY --from=builder /app/config/url_encoder/config-docker.yaml /app/config/url_encoder/config-docker.yaml
COPY --from=builder /app/url_encoder /app/url_encoder
ENTRYPOINT ["/app/url_encoder"]
