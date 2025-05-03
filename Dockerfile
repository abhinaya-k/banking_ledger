# ---------- Build Stage ----------
FROM golang:alpine as builder

WORKDIR /app

# Install build tools and migration dependencies
RUN apk update && apk add --no-cache \
    gcc \
    libc-dev \
    librdkafka-dev \
    pkgconf \
    curl \
    tar

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire app source
COPY . .

COPY database /database


# Install migrate binary (for Alpine/Linux)
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
    -o migrate.tar.gz && \
    tar -xvf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    rm migrate.tar.gz

# Build the Go binary
RUN go build -tags musl -o goserver

# ---------- Final Stage ----------
FROM alpine:latest

WORKDIR /app

# Install PostgreSQL client (for pg_isready)
RUN apk add --no-cache postgresql-client

# Copy Go server binary and migration tool from builder
COPY --from=builder /app/goserver /app
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/database /app/database

# Expose service port
EXPOSE 8080

# Let docker-compose override this with its custom entrypoint (that waits, runs migration, and starts the service)
ENTRYPOINT ["./goserver"]
