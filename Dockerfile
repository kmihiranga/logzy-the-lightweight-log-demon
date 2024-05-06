# use an official Go runtime as a parent image
FROM golang:1.21-alpine3.19 as builder

# Set the working directory inside the container
WORKDIR /app

# Create a directory for log files or other purposes
RUN mkdir /

# Download Go modules
COPY go.mod go.sum ./
COPY ops/ ./
RUN go mod download

# Install necessary packages (optional)
# For example, if you need additional utilities:
RUN apk add --no-cache bash

# Create a directory within /var/log for the log
RUN mkdir -p /var/log/logzy

# Set permissions for /var/log/logzy
RUN chmod 755 /var/log/logzy

# Create a log file for the application logging
RUN touch /var/log/logzy/logzy.log

# Set the permission to the created file
RUN chmod 644 /var/log/logzy/logzy.log

# Copy the local package files to the container's workspace
COPY . .

# Build your application on the specified directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logzy .

# Use a docker multi-stage build to minimize the size of the final image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/logzy .

# Create a ops directory
RUN mkdir -p ./ops

COPY --from=builder /app/ops ./ops

# Optional - not required for the production environment - for testing purposes
COPY --from=builder /app/netwiz_account_service_log.log /var/log/netwiz_account_service_log.log

# Make sure the logzy binary is executable
RUN chmod +x ./logzy

# Optionally set permissions for the log file - for testing purposes
RUN chmod 644 /var/log/netwiz_account_service_log.log

# optionally create a directoryu for log
RUN mkdir -p /var/log/logzy

# Run the application when the container starts
CMD [ "./logzy" ]