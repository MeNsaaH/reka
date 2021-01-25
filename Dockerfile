# Builder image so to add ca-certificates to scratch
FROM golang:alpine as build
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY reka  /
ENTRYPOINT ["/reka"]