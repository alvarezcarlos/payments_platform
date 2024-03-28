# Use a minimal and secure base image
FROM golang:1.22.1 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application directory into the container
COPY app .

# Build the Go app with necessary flags for security and optimization
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app .

# Start a new stage from a lightweight base image
FROM alpine:latest

# Set environment variables for configuration
ENV PORT=8080
ENV DB_HOST=postgres
ENV DB_NAME=payments
ENV DB_PASSWORD=secret123
ENV DB_PORT=5432
ENV DB_USERNAME=postgres
ENV ENV=local

# Copy the pre-built binary from the previous stage
COPY --from=build /app/app /app/app

# Set permissions for the binary
RUN chmod +x /app/app

# Expose the port on which the application will listen
EXPOSE $PORT

# Run the application
CMD ["/app/app"]

