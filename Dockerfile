FROM golang:alpine AS builder

# Need git for dependencies
RUN apk update && apk add --no-cache git

# Create working dir
WORKDIR /build

# Copy source
COPY . .

# Get dependencies
RUN go get -d -v

# Build app
RUN go build -o /app/bin/main

# Copy config
COPY config/example.config.docker.yml /app/config/config.yml

FROM alpine

# Copy only config and executable
COPY --from=builder /app/config /config
COPY --from=builder /app/bin /app/

EXPOSE 9000

# Run app
CMD ["./app/main"]