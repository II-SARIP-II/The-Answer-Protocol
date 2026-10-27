package server

import (
	"log/slog"
	"bytes"
	"os"
)

type logWriter struct{}

func SetLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(logWriter{}, nil)))
}

func (w logWriter) Write(p []byte) (int, error) {
	if bytes.Contains(p, []byte(`"WARN"`)) || bytes.Contains(p, []byte(`"ERROR"`)) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

