package services

import "log/slog"

type Service interface {
	GetLogger() *slog.Logger
	GetConnection() any
}
