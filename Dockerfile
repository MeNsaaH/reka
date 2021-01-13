FROM golang:1.15 AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github/mensaah/reka/
COPY . .
RUN go get -v
RUN go build -o /go/bin/reka

FROM scratch
COPY --from=builder /go/bin/reka  /go/bin/reka
ENTRYPOINT ["/go/bin/reka"]