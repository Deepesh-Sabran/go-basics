package helper

import (
	"errors"

	"github.com/Deepesh-Sabran/go-basics/internal/models"
)

func SimulateSendEmail(job models.EmailJob) error {

	if job.Email == "" {
		return errors.New("invalid recipient")
	}

	if job.UserID%2 == 0 {
		return errors.New("temporary email provider timeout")
	}

	return nil
}