FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/server ./services/products-service


# Add these environment variables to fix proxy issues
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on


FROM alpine:latest
RUN apk update && apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 50052
CMD ["./server"]