FROM golang:1.26.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/loadsg \
    ./cmd/server

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/loadsg-cli \
    ./cmd/loadsg


FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/loadsg ./loadsg
COPY --from=builder /out/loadsg-cli ./loadsg-cli
COPY config.json ./config.json

EXPOSE 8080

ENTRYPOINT ["./loadsg"]
