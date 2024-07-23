FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go get -u github.com/swaggo/swag/cmd/swag

# Build the Go application
RUN go build -o main .

EXPOSE ${WEBSERVER_PORT}
CMD ["./main"]
