# Build Backend
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app cmd/main.go

# Final Image
FROM debian:bookworm-slim
WORKDIR /app

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy backend binary
COPY --from=builder /app/bin/app ./app

# Copy migrations
COPY migrations ./migrations

# Create uploads directory
RUN mkdir -p uploads

# Default environment variables
ENV APP_PORT=8765
ENV APP_ENV=production
ENV UPLOAD_DIR=/app/uploads

# Expose the application port
EXPOSE 8765

# Run the application
CMD ["./app"]
