# Stage 1: Builder
# Use a Go base image with the required Go version for building
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container for the build process
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's build cache
# If these files don't change, Docker won't re-download modules
COPY go.mod go.sum ./

ENV GOPROXY "https://goproxy.cn,direct"

# Download all Go modules
# CGO_ENABLED=0 is important for creating a statically linked binary
# This means the binary won't depend on system libraries and can run on minimal base images
RUN go mod download

# Copy the rest of your application source code
COPY . .

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

# Build the Go application
# -o main: specifies the output binary name
# -ldflags="-w -s": strips debug information and symbol tables, further reducing binary size
# GOOS=linux: ensures the binary is compiled for Linux, regardless of the host OS
# GOARCH=amd64: ensures the binary is compiled for amd64 architecture
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o fileshare .

# Stage 2: Runner
# Use a minimal base image like scratch or alpine for the final image
# scratch is the smallest possible image, containing nothing but your binary
# alpine is slightly larger but includes basic utilities if you need them (e.g., for debugging)
FROM alpine:latest

# Install ca-certificates for secure connections (HTTPS) if your app makes external calls
# This is crucial if your Go app communicates with other services over TLS
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache sqlite-libs


# Set the working directory for the final application
WORKDIR /app

# Copy the built binary from the 'builder' stage to the final image
COPY --from=builder /app/fileshare .

COPY settings.yml .
# Expose the port your Go application listens on
EXPOSE 60011

EXPOSE 8080

# Define the command to run your application
CMD ["/app/fileshare", "server"]
