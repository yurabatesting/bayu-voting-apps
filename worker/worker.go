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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("voting-worker"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	ctx := context.Background()

	// Inisialisasi OpenTelemetry Tracer Provider
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error saat mematikan tracer provider: %v", err)
		}
	}()

	tracer := otel.Tracer("voting-worker-tracer")

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis:6379"
	} else {
		redisHost = redisHost + ":6379"
	}

	queueName := os.Getenv("REDIS_QUEUE")
	if queueName == "" {
		queueName = "votes"
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "postgres"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	rdb := redis.NewClient(&redis.Options{Addr: redisHost})
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Worker berjalan. Menunggu antrean dari Redis...")

	for {
		result, err := rdb.BLPop(ctx, 0, queueName).Result()
		if err == nil {
			vote := result[1]

			// Membuat span tracing untuk proses kerja worker
			_, span := tracer.Start(ctx, "ProcessVote", trace.WithSpanKind(trace.SpanKindInternal))

			_, dbErr := db.Exec("INSERT INTO votes (vote) VALUES ($1)", vote)
			if dbErr != nil {
				log.Println("Gagal menyimpan ke DB:", dbErr)
			} else {
				log.Println("Berhasil memproses vote:", vote)
			}

			span.End()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
