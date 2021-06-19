FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o main .

WORKDIR /app

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /build/main .

EXPOSE 8080

CMD ["./main"]
