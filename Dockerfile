# Dockerfile References: https://docs.docker.com/engine/reference/builder/
#
# Cert stage
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# Builder stage
FROM golang:1.13.4-alpine3.10 as builder

# Output dir
RUN mkdir -p /build

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy mod file inside the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source inside the container
COPY . .

# Compile output
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -installsuffix cgo -o /bin/dotheting cmd/dothething/main.go

# Thin stage
FROM alpine:3.11.3
RUN apk add --no-cache ca-certificates openssl

ENV PATH=/bin

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/dothething /bin/dothething

ENTRYPOINT ["/bin/dothething"]

