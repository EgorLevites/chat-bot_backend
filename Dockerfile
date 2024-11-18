# Build stage using an image based on Debian 12
FROM golang:1.21-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the project source code
COPY . .

# Build the application
RUN go build -o main .

# Final image using distroless based on Debian 12
FROM gcr.io/distroless/base-debian12

# Copy the compiled binary from the build stage
COPY --from=builder /app/main /

# Expose the port the application will use
EXPOSE 8080

# Specify the command to run the application
CMD ["/main"]
