# Use the official Golang image as a base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git to fetch Go modules
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application for ARM64 architecture on Linux OS
RUN GOARCH=arm64 GOOS=linux go build -o myCode .

# Expose port 8088 to the outside world
EXPOSE 8088

# Command to run the executable
CMD ["./myCode"]
