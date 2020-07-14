

# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.12.8 as builder

# Add Maintainer Info
LABEL maintainer="Nuryanto <nuryantofattih@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod , sum files config.json
COPY go.mod go.sum config.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/


# Build Args
ARG LOG_DIR=/root/wwwroot


# Create Log Directorytail
RUN mkdir -p ${LOG_DIR}


# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

COPY --from=builder /app/config.json .

# Expose port 8080 to the outside world
EXPOSE 9191

# Declare volumes to mount
VOLUME [${LOG_DIR}]

# Command to run the executable
CMD ["./main"] 
