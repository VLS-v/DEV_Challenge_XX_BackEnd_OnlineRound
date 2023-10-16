FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

# Build the Go app
RUN go build -o spreadsheets main.go

FROM alpine
WORKDIR /bin

COPY --from=builder /app/spreadsheets /bin/spreadsheets
COPY --from=builder /app/saves /bin/saves

EXPOSE 8080

CMD ["spreadsheets"]
