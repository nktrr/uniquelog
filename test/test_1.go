package test

import (
	"log/slog"
)

func a(lol *slog.Logger) {
	lol.Info("error1")
	lol.Warn("error2")
	lol.Info("error1")
}

type Mock struct {
	log *slog.Logger
}

func (m Mock) c() {
	m.log.Info("error1")
}
