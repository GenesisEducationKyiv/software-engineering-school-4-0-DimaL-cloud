FROM golang:1.22-alpine
WORKDIR /rate-service
COPY .. .
RUN go build -o main ./cmd/main.go
EXPOSE 8080
ENTRYPOINT ["/rate-service/main"]