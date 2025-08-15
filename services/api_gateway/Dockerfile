########################################
# Builder stage
########################################
FROM golang:1.24-alpine AS builder

# Optional: set a build-time arg for binary name and module download mode
ARG BINARY_NAME=app
ARG CGO_ENABLED=0
ARG GOARCH=amd64
ARG GOOS=linux

# Set environment variables from arg
ENV CGO_ENABLED=$CGO_ENABLED
ENV GOOS=$GOOS
ENV GOARCH=$GOARCH

WORKDIR /app

# Copy go mod first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
RUN ls

COPY . .

RUN go build -o ${BINARY_NAME} ./cmd

########################################
# Runtime stage (Alpine)
########################################
FROM alpine:3.18

# Add CA certificates for HTTPS if your app makes outbound TLS calls
RUN apk add --no-cache ca-certificates

# Create a non-root user to run the app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/${BINARY_NAME} /app/${BINARY_NAME}

RUN chmod +x /app/${BINARY_NAME}

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/app"]