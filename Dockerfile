FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /spotify-car

FROM alpine:latest

WORKDIR /

COPY --from=builder /spotify-car /spotify-car

EXPOSE 8080

ENTRYPOINT ["/spotify-car"]
