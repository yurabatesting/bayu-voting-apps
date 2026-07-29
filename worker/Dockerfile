FROM golang:1.22-alpine
WORKDIR /app
COPY worker.go .
# Inisialisasi module dan download library yang dibutuhkan
RUN go mod init worker && go get github.com/go-redis/redis/v8 github.com/lib/pq
RUN go build -o worker worker.go
CMD ["./worker"]