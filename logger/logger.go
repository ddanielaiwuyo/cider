package logger

import (
	"log/slog"
	"os"
)

func JSONLogger() *slog.Logger {
	var out *os.File
	var err error
	out, err = os.OpenFile("logs/logs.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		slog.Warn("could not find log/loggs.txt, using stdout")
		out = os.Stdout
		slog.Error("Err:", "", err)
	}

	logger := slog.New(slog.NewJSONHandler(out, nil))
	return logger
}
