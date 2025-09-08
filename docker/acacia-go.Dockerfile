FROM golang:1.22-alpine

WORKDIR /app

RUN addgroup -g 1001 -S golang && \
    adduser -S golang -u 1001

# Copy go mod files first for better caching
COPY services/acacia-go/go.mod services/acacia-go/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY services/acacia-go/ ./

RUN chown -R golang:golang /app
USER golang

EXPOSE 8080

CMD ["go", "run", "main.go"]