package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/Deepesh-Sabran/go-basics/internal/helper"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
	"github.com/redis/go-redis/v9"
)

func StartEmailWorker(ctx context.Context) {
	for {
		select {
		case <- ctx.Done():
			log.Println("email worker stopped")
			return
		default:
		}

		result, err:= config.RedisClient.LMove(
			config.Ctx,
			"email_jobs",
			"email_processing",
			"RIGHT",
			"LEFT",
		).Result()

		if err == redis.Nil {
			time.Sleep(2 * time.Second)
			continue
		}

		if err != nil {
			log.Println("Worker error:", err)
			continue
		}

		var job models.EmailJob

		if err:= json.Unmarshal([]byte(result), &job); err != nil {
			log.Println("Invalid Job", err)

			config.RedisClient.LRem(
				config.Ctx,
				"email_processing",
				1,
				result,
			)

			continue
		}

		sendErr:= helper.SimulateSendEmail(job)

		if sendErr != nil {
			if job.Retries < 3 {
				job.Retries++

				data, _:= json.Marshal(job)

				config.RedisClient.LPush(
					config.Ctx,
					"email_jobs",
					data,
				)

				config.RedisClient.LRem(
					config.Ctx,
					"email_processing",
					1,
					result,
				)

				log.Printf("retrying email job (attempt %d)", job.Retries)
			} else {
				config.RedisClient.LPush(
					config.Ctx,
					"failed_email_jobs",
					result,
				)

				log.Println("moved job to failed queue")
			}

			continue
		}

		log.Printf(
			"📨 Sending welcome email to Id: %d and Name: %s",
			job.UserID,
			job.Name,
		)

		config.RedisClient.LRem(
			config.Ctx,
			"email_processing",
			1,
			result,
		)
	}
}