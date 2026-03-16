## Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bot ./main.go

## Runtime stage
FROM alpine:3.23

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /bot /app/bot

# TOKEN и USERNAME должны приходить через окружение или --env-file
ENV TOKEN=""
ENV USERNAME=""

CMD ["/app/bot"]

