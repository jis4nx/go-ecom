# Builder Stage
FROM golang:1.22.2-bookworm as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o bin/product product/cmd/main.go


# Final Stage
FROM alpine:3.19

COPY --from=builder /app/bin/product .

CMD ["./product"]

