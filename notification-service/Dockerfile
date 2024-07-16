FROM golang:1.22-alpine
WORKDIR /notification-service
COPY .. .
RUN go build -o main ./cmd/main.go
EXPOSE 8080
ENTRYPOINT ["/notification-service/main"]