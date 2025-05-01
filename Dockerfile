FROM golang:1.24-alpine
LABEL org.opencontainers.image.source="https://github.com/Rherer/restic-exporter"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Install restic
RUN apk add restic

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /restic-exporter

# Just for documentation, port will still need to be exposed via docker run
EXPOSE 8080

# Run
CMD ["/restic-exporter"]