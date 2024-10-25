FROM golang:1.23.0 AS builder
WORKDIR /module
COPY go.mod go.sum  ./
RUN go mod download
COPY .  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:3.20.3
RUN apk --no-cache add ca-certificates && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    echo "Asia/Bangkok" >  /etc/timezone && \
    apk del tzdata
WORKDIR /root/
COPY --from=builder /module/app .
ENV GIN_MODE release
EXPOSE 8000
CMD ["./app"]