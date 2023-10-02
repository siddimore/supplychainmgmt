# Use an official Golang runtime as a parent image
FROM golang:1.17

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Build the Go application
RUN go get -d -v ./...
RUN go install -v ./...

# Expose port 8080 for the application
EXPOSE 8080

# Run the application when the container starts
CMD ["app"]
