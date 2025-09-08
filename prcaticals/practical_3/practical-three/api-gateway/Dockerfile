# api-gateway/Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/server ./api-gateway

# Add these environment variables to fix proxy issues
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on


FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]