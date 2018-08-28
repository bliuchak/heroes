FROM golang as builder

COPY cmd /go/src/github.com/bliuchak/heroes/cmd
COPY internal /go/src/github.com/bliuchak/heroes/internal
COPY Gopkg.lock /go/src/github.com/bliuchak/heroes
COPY Gopkg.toml /go/src/github.com/bliuchak/heroes

WORKDIR /go/src/github.com/bliuchak/heroes

RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo cmd/heroes/main.go
RUN ls -la && pwd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/bliuchak/heroes/main .
CMD ["./main"]
