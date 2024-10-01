# Start from the official Golang base image for building the app
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image to run the application
FROM scratch

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside
EXPOSE 8080

# Run the compiled Go binary
CMD ["./main"]
