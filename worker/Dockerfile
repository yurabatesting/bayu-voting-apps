FROM golang:1.22-alpine

WORKDIR /app

COPY worker.go .

# Inisialisasi module dan download seluruh library (termasuk OpenTelemetry)
RUN go mod init worker && \
    go get github.com/go-redis/redis/v8 && \
    go get github.com/lib/pq && \
    go get go.opentelemetry.io/otel && \
    go get go.opentelemetry.io/otel/sdk && \
    go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp && \
    go get go.opentelemetry.io/otel/trace

RUN go build -o worker worker.go

CMD ["./worker"]