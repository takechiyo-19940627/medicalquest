FROM golang:1.21-alpine

WORKDIR /app

# Install development tools
RUN apk add --no-cache git make gcc libc-dev

# Install air for hot reloading in development
RUN go install github.com/cosmtrek/air@latest

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Create .air.toml if it doesn't exist
RUN if [ ! -f .air.toml ]; then \
    echo 'root = "."' > .air.toml && \
    echo 'tmp_dir = "tmp"' >> .air.toml && \
    echo '[build]' >> .air.toml && \
    echo '  cmd = "go build -o ./tmp/main ."' >> .air.toml && \
    echo '  bin = "./tmp/main"' >> .air.toml && \
    echo '  delay = 1000' >> .air.toml && \
    echo '  exclude_dir = ["assets", "tmp", "vendor"]' >> .air.toml && \
    echo '  include_ext = ["go", "tpl", "tmpl", "html"]' >> .air.toml && \
    echo '  exclude_regex = ["_test\\.go"]' >> .air.toml; \
    fi

EXPOSE 8080

# Use air for live reloading in development
CMD ["air", "-c", ".air.toml"]