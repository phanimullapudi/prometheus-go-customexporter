FROM golang:alpine3.16 as builder
RUN mkdir /src
ADD . /src/
WORKDIR /src
RUN go build -o iscalecc-prometheus-exporter
FROM alpine
COPY --from=builder /src/iscalecc-prometheus-exporter /app/iscalecc-prometheus-exporter
WORKDIR /app
ENTRYPOINT ["/app/iscalecc-prometheus-exporter"]