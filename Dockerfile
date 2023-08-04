# Start from the Go builder image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code to the container
COPY . .

# Run 'go mod tidy' to update dependencies
RUN go mod tidy

# Build the Go program inside the container
RUN go build -o /myapp-operator

# Set the binary as the entrypoint
ENTRYPOINT ["/myapp-operator"]

