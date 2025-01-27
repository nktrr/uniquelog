package test

import (
	"log/slog"
)

func a(lol *slog.Logger) {
	lol.Info("error1")
	lol.Warn("error2")
	lol.Info("error3")
}
