FROM golang:1.17.2 AS builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o oopsie main.go

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /workspace/oopsie .
CMD ["/oopsie"]
