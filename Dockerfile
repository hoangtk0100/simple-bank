# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db.env .
COPY start.sh .
RUN chmod +x start.sh
COPY wait-for.sh .
RUN chmod +x wait-for.sh
COPY db/migration ./db/migration

EXPOSE 8088
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]