# Builder Stage
FROM golang:1.22.2-bookworm as builder

WORKDIR /app

COPY product/ ./product
COPY helpers/ ./helpers
COPY config/ ./config

COPY go.mod .
COPY go.sum .

COPY dev.env .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/product product/cmd/main.go


# Final Stage
FROM alpine:3.19

COPY --from=builder /app/bin/product .
COPY --from=builder /app/dev.env .

CMD ["./product"]
