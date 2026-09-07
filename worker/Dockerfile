FROM golang:1.20-alpine

WORKDIR /app

COPY worker.go .

# Inisialisasi modul dan tentukan versi stabil secara eksplisit
RUN go mod init worker && \
    go mod edit -require=github.com/go-redis/redis/v8@v8.11.5 && \
    go mod edit -require=github.com/lib/pq@v1.10.9 && \
    go mod edit -require=go.opentelemetry.io/otel@v1.24.0 && \
    go mod edit -require=go.opentelemetry.io/otel/sdk@v1.24.0 && \
    go mod edit -require=go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp@v1.24.0 && \
    go mod edit -require=go.opentelemetry.io/otel/trace@v1.24.0 && \
    go mod tidy

RUN go build -o worker worker.go

CMD ["./worker"]