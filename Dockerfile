# Use a smaller base image for the build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /src
COPY . .

# Enable Go modules
ENV GO111MODULE=on
RUN go mod tidy
# Build the Go app and strip debugging information to reduce size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/action ./

FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the compiled Go program from the builder stage
COPY --from=builder /bin/action /app/action

# Set entrypoint
ENTRYPOINT ["/app/action", "run"]