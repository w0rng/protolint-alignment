FROM golang:1.24-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o plugin main.go &&  \
    go install github.com/yoheimuta/protolint/cmd/protolint@latest

FROM alpine:latest
WORKDIR /app
LABEL io.whalebrew.name=protolint
COPY --from=builder /build/plugin .
COPY --from=builder /go/bin/protolint .
ENTRYPOINT ["/app/protolint", "-plugin", "/app/plugin"]
