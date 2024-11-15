FROM golang:1.23-alpine3.20 AS builder
RUN apk add make
WORKDIR /go/app
COPY . .
RUN make build

FROM alpine:3.20
COPY --from=builder /go/app/.dist/fake-dcgm-exporter /usr/local/bin/fake-dcgm-exporter
ENTRYPOINT [ "fake-dcgm-exporter" ]
