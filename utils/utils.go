package utils

import (
	"JeopardyScoreBoardV2/models"
)

// Helper function to allow repos to send result payloads with less text.
func GetResult[T any](err error, statusCode int, payload T) models.Result[T] {
	return models.Result[T]{
		StatusCode: statusCode,
		Err:        err,
		ResultData: payload,
	}
}
