package workers

import (
	"log"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/redis/go-redis/v9"
)

func RecoverEmailJobs() {
	for {
		job, err:= config.RedisClient.RPopLPush(config.Ctx, "email_processing", "email_jobs").Result()

		if err != redis.Nil {
			break
		}

		if err != nil {
			log.Println("Email recovery error: ", err)
			break
		}

		log.Println("Recovered stuck job", job)
	}
}