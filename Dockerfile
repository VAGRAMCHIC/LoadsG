FROM golang:1.23.10-alpine AS build

RUN apk add --no-cache git
WORKDIR /app

RUN go mod init loadsg
RUN go mod tidy
RUN go get github.com/gin-gonic/gin
RUN go get github.com/golang-jwt/jwt/v5
RUN go get github.com/jackc/pgx/v5
COPY . .

RUN go build -o app .


FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=build /app/app .
COPY config.json /app/
EXPOSE 8080
CMD ["./app"]
