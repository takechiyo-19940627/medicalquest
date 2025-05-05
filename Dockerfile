FROM golang:1.23-alpine

WORKDIR /app

# Install development tools
RUN apk add --no-cache git make gcc libc-dev

# Install air for hot reloading in development
RUN go install github.com/air-verse/air@latest


# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the application
COPY . .

EXPOSE 8080

# Use air for live reloading in development
CMD ["air", "-c", ".air.toml"]