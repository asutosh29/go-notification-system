FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd 

FROM alpine:3.21
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates


EXPOSE 8080
CMD ["./main"]
