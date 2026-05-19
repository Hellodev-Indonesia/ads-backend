package main

import (
	"log"

	"github.com/alex/ads_backend/config"
	"github.com/hibiken/asynq"
)

func main() {
	config.LoadEnv()
	
	redisHost := config.GetEnv("REDIS_HOST", "localhost")
	redisPort := config.GetEnv("REDIS_PORT", "6379")
	
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	// mux.HandleFunc(tasks.TypeMetaSync, tasks.HandleMetaSyncTask)

	log.Println("Worker starting...")
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
