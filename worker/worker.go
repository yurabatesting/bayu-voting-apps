package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
    _ "github.com/lib/pq"
)

func main() {
    ctx := context.Background()

    // Mengambil konfigurasi dari Environment Variables
    redisHost := os.Getenv("REDIS_HOST")
    if redisHost == "" {
        redisHost = "redis:6379"
    } else {
        redisHost = redisHost + ":6379"
    }

    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbHost := os.Getenv("DB_HOST")

    // Membuat string koneksi database secara dinamis
    dsn := fmt.Sprintf("postgres://%s:%s@%s/postgres?sslmode=disable", dbUser, dbPassword, dbHost)

    // Konek ke Redis
    rdb := redis.NewClient(&redis.Options{Addr: redisHost})

    // Konek ke PostgreSQL
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Worker berjalan. Menunggu antrean dari Redis...")

    for {
        result, err := rdb.BLPop(ctx, 0, "votes").Result()
        if err == nil {
            vote := result[1]
            _, err = db.Exec("INSERT INTO votes (vote) VALUES ($1)", vote)
            if err != nil { 
                log.Println("Gagal menyimpan ke DB:", err) 
            } else {
                log.Println("Berhasil memproses vote:", vote)
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}