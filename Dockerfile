FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build the binary with a static build for a smaller size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /docker-gs-ping

# Include the configuration file (assuming it's named config.yaml)
COPY ./config.yaml /config.yaml

# Final stage
FROM alpine:latest

# Copy the binary from the builder image
COPY --from=builder /docker-gs-ping /docker-gs-ping

# Copy the configuration file
COPY --from=builder /config.yaml /config.yaml

# Expose port 9000 (adjust this to match your application's config)
EXPOSE 9000

# Define an entry point
ENTRYPOINT ["/docker-gs-ping"]

# Specify default command-line arguments (if needed)
CMD ["us=address", "ru=address2"]
