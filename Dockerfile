# Use the official Golang image as a build stage
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o main .

# Use a smaller base image for the final container
FROM gcr.io/distroless/base-debian11

# Set working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Command to run the binary
ENTRYPOINT ["/main"]
