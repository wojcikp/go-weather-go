FROM golang:1.22.1-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o weather-hub ./cmd/app
EXPOSE 8081
CMD ["./weather-hub"]
