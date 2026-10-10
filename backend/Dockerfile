# Step 1: Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependency downloads
COPY go.mod ./
# If go.sum exists: COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments passed from Jenkins / docker build
ARG VERSION=1.0.0
ARG GIT_COMMIT=unspecified
ARG BUILD_TIME=unspecified

# Compile statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'" \
    -o /app/server .

# Step 2: Final runtime stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]
