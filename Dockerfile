FROM golang:1.23.10-alpine AS build

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o app .


FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=build /app/app .
COPY config.json /app/
EXPOSE 8080
CMD ["./app"]
