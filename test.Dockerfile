FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

# Build the Go app
RUN go test -c -o spreadsheets.test ./controllers/

FROM alpine
WORKDIR /bin

COPY --from=builder /app/spreadsheets.test /bin/spreadsheets.test
COPY --from=builder /app/saves /bin/saves

EXPOSE 8080

CMD ["spreadsheets.test", "-test.v"]
