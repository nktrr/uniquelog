package test

import (
	"log/slog"
)

type Mock struct {
	log *slog.Logger
}

func (m Mock) c() {
	m.log.Info("error1")
}
