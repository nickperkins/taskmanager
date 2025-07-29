# ----------- Builder Stage -----------
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o taskmanager ./cmd/server

# ----------- Runtime Stage -----------
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/taskmanager ./taskmanager
EXPOSE 8080
ENTRYPOINT ["/app/taskmanager"]
