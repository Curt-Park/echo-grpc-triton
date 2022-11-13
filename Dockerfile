FROM golang:alpine3.16 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build
COPY . .
RUN go build *.go

WORKDIR /dist
RUN cp /build/grpc_service.pb .

FROM scratch
COPY --from=builder /dist/grpc_service.pb .
ENTRYPOINT ["/grpc_service.pb"]
