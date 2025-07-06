# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary and config files
COPY config.yaml /app/
COPY bin/turbogin-linux /app/

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./turbogin-linux"]