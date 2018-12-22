FROM golang:1.11.4 as builder

COPY . /opt/heroes
WORKDIR /opt/heroes

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo cmd/heroes/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /opt/heroes .
CMD ["./main"]
