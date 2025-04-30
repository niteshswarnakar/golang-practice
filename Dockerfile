FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod tidy

COPY . ./

# Build the Go app
RUN go build -o main .

# STAGE 2: CREATE A MINIMAL RUNTIME IMAGE
FROM debian:bookworm

WORKDIR /app

# Copy the binary built in the first stage
COPY --from=builder /app .

EXPOSE 5000

# Run the compiled binary directly, not with "go run"
CMD ["./main"]