package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
}

func (processor *RedisTaskProcessor) ProcessVerifyEmailTask(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// TODO

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", "#").Msg("processed task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessCreateUserTask(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	return nil
}
