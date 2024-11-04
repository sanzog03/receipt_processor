# Use official Go runtime
FROM golang:1.23.2 as builder

# Ensure Go compiler builds a static linked binary
ENV CGO_ENABLED=0

WORKDIR /receiptProcessor
 
COPY go.mod  ./
RUN go mod download

# Install make and git
RUN apt update && apt install make git

# Copy the application source code
COPY . .

# Install swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
# Build application
RUN make all

# new stage
FROM alpine:latest 

RUN mkdir /configs
COPY --from=builder /receiptProcessor/main /usr/local/bin/
COPY --from=builder /receiptProcessor/configs/config.yml /configs/config.yml
# Give Execution permission
RUN chmod +x /usr/local/bin/main

# Expose port 9080
EXPOSE 9080

# Entrypoint
CMD ["main"]
