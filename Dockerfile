### BASE - USED IN ALL STEPS
################################################################################
FROM golang:1.22.2-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
ENV GOOS=linux

### BUILD DEPENDENCIES
################################################################################
FROM base AS build-dependencies
COPY go.mod go.sum ./
RUN go mod download

### BUILDER
################################################################################
FROM build-dependencies AS builder
COPY . .
RUN go build -ldflags="-w -s" -o /app/ics-generator-service ./cmd/server

### PRODUCTION SERVER
################################################################################
FROM alpine:3.19 AS production
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/ics-generator-service .

# Expose the port the app runs on
EXPOSE 8080

# Use a non-root user
RUN adduser -D appuser
USER appuser

# Command to run the executable
ENTRYPOINT ["./ics-generator-service"]
CMD ["--server"]