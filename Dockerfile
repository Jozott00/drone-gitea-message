# Use the official Golang image from Docker Hub:
FROM golang:latest as builder

# Set the working directory inside the container:
WORKDIR /app

# Copy the local code to the container's workspace:
COPY . .

# Build the Go app for the amd64 architecture:
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 GO111MODULE=on go build -o /drone-gitea-message .

# Start a new stage for a smaller final image:
FROM plugins/base:multiarch

# Copy the binary from the builder stage:
COPY --from=builder /drone-gitea-message /bin/

# Run the binary:
ENTRYPOINT ["/bin/drone-gitea-message"]
