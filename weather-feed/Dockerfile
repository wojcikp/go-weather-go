FROM golang:1.22.1-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS build
COPY . .
RUN go build -o weather-feed ./cmd/app

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/weather-feed .
COPY --from=build /app/assets /app/assets
COPY --from=build /app/config /app/config
CMD ["./weather-feed"]
