FROM golang:1.22.2-bookworm

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY services/ ./services
COPY helpers/ ./helpers
COPY config/ ./config
COPY pkg/ ./pkg


COPY dev.env .

RUN CGO_ENABLED=0 GOOS=linux

RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

RUN export PATH=$PATH:$(go env GOPATH)/bin
CMD ["air", "-c", "./services/product/.air.toml"]
