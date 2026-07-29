package main

import (
    "context"
    "database/sql"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
    _ "github.com/lib/pq"
)

func main() {
    ctx := context.Background()
    // Konek ke Redis
    rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})

    // Konek ke PostgreSQL
    db, err := sql.Open("postgres", "postgres://postgres:password@db/postgres?sslmode=disable")
    if err != nil { log.Fatal(err) }

    log.Println("Worker berjalan. Menunggu antrean dari Redis...")

    for {
        // Mengambil antrean 'votes' dari Redis
        result, err := rdb.BLPop(ctx, 0, "votes").Result()
        if err == nil {
            vote := result[1]
            // Menyimpan ke PostgreSQL
            _, err = db.Exec("INSERT INTO votes (vote) VALUES ($1)", vote)
            if err != nil { 
                log.Println("Gagal menyimpan ke DB:", err) 
            } else {
                log.Println("Berhasil memproses vote:", vote)
            }
        }
        time.Sleep(100 * time.Millisecond) // Mencegah CPU spike
    }
}