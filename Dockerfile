FROM golang:1.13.4
WORKDIR /go/src/github.com/afritzler/oopsie
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/afritzler/oopsie/oopsie .
CMD ["/oopsie"]