# First stage: Build the Go application
FROM golang:1.21-alpine as build

# Install make and other necessary tools
RUN apk add --no-cache make bash gcc musl-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download
RUN go mod vendor

# Install oapi-codegen for generating OpenAPI 3 code
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Ensure Go bin directory is in the PATH
ENV PATH=$PATH:/go/bin

# Copy the source code into the container
COPY . .

# Generate OpenAPI 3 code
RUN make generated

# Build the Go app with CGO disabled for static linking
RUN CGO_ENABLED=0 go build -o main ./cmd

# Second stage: Create a smaller image for production
FROM alpine:latest

# Install necessary utilities
RUN apk add --no-cache bash

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
