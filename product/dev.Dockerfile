FROM golang:1.22.2-bookworm

WORKDIR /app

COPY product/ ./product
COPY helpers/ ./helpers
COPY config/ ./config

COPY go.mod .
COPY go.sum .

COPY dev.env .

RUN CGO_ENABLED=0 GOOS=linux
RUN go mod download
RUN go install github.com/cosmtrek/air@latest
RUN export PATH=$PATH:$(go env GOPATH)/bin
CMD ["air", "-c", "./product/.air.toml"]
