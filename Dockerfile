FROM golang:1.17.1 AS builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o oopsie main.go

FROM alpine:3.14.2
RUN apk --no-cache add ca-certificates=20191127-r5
WORKDIR /
COPY --from=builder /workspace/oopsie .
CMD ["/oopsie"]
