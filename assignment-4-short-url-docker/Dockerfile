# Stage 1: Build stage
FROM golang:1.22.4-alpine AS build

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o short-url .

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/short-url .

EXPOSE 8080

# Set the entrypoint command
CMD ["./short-url"]