package test

import "log/slog"

func b(log *slog.Logger) {
	log.Info("error1")
}
