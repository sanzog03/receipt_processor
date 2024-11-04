# Use official Go runtime
FROM golang:1.21.1 as builder

# Ensure Go compiler builds a static linked binary
ENV CGO_ENABLED=0

WORKDIR /receiptProcessor
 
COPY go.mod  ./
RUN go mod download

# Install make and git
RUN apt update && apt install make git

# Copy the application source code
COPY . .
RUN ls -l

# Build application
RUN make all

# Staging 
FROM alpine:latest 

RUN mkdir /config
COPY --from=builder /receiptProcessor/main /usr/local/bin/
COPY --from=builder /receiptProcessor/config/config.yml /config/config.yml
# Give Execution permission
RUN chmod +x /usr/local/bin/main

# Expose port 808
EXPOSE  8080

# Entrypoint
CMD ["main"]
