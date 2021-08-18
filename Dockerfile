FROM golang:1.17.0 AS builder
WORKDIR /go/src/github.com/afritzler/oopsie
COPY . .
ENV GO111MODULE=on
RUN make build

FROM alpine:3.14.1
RUN apk --no-cache add ca-certificates=20191127-r5
WORKDIR /
COPY --from=builder /go/src/github.com/afritzler/oopsie/oopsie .
CMD ["/oopsie"]
