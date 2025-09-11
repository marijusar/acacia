FROM golang:1.24-alpine

ARG USER_ID=1000
ARG GROUP_ID=1000

WORKDIR /app

# Create user with host user/group IDs
RUN addgroup -g ${GROUP_ID} -S appgroup && \
    adduser -S appuser -u ${USER_ID} -G appgroup

# Set up Go module cache directory with proper permissions
RUN mkdir -p /go/pkg/mod && \
    chown -R ${USER_ID}:${GROUP_ID} /go

# Copy go mod files first for better caching
COPY services/acacia/go.mod services/acacia/go.sum ./
RUN chown ${USER_ID}:${GROUP_ID} go.mod go.sum

# Download dependencies as the app user
USER appuser
RUN go mod download

# Copy source code
COPY services/acacia/ ./

EXPOSE 8080

CMD ["go", "run", "main.go"]
