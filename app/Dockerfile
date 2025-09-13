# Multi-stage build for Go server
FROM golang:1.22 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server

# Minimal runtime image
FROM scratch

ENV PORT=5000
EXPOSE 5000

COPY --from=builder /server /server

ENTRYPOINT ["/server"]
