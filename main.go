package main

import (
	"log/slog"
	"os"
	"secure-chat-client/env"
	"secure-chat-client/gui"
)

func initSlog() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))
}

func main() {
	env.LoadEnv()
	initSlog()
	gui.Start()
}
