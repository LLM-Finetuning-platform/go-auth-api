FROM golang:1.26.4-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o auth-grpc ./cmd

# Stage to execute the binary

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/auth-grpc /app/auth-grpc

EXPOSE 50051

USER nonroot:nonroot

ENTRYPOINT ["/app/auth-grpc"]