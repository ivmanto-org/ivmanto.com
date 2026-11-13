# Stage 1: Build the Go binary in a dedicated build environment.
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker layer caching for dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the application into a static binary.
# CGO_ENABLED=0 is crucial for a static binary in Alpine.
# -ldflags="-w -s" strips debug information to create a smaller binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o ivmanto-backend-service ./cmd/server

# booking-cleanup is the booking-data retention job (Cloud Run Job,
# scheduled daily via Cloud Scheduler) — it ships in the same image so
# the backend image serves both the service and the job.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o booking-cleanup ./cmd/booking-cleanup

# Stage 2: Create the final, minimal image using distroless.
FROM gcr.io/distroless/static-debian11

# Copy only the compiled binaries from the builder stage.
COPY --from=builder /app/ivmanto-backend-service /
COPY --from=builder /app/booking-cleanup /booking-cleanup

# Set the entrypoint for the container. Cloud Run will provide the PORT.
ENTRYPOINT ["/ivmanto-backend-service"]
