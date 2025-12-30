package main

import (
	"secure-chat-client/env"
	"secure-chat-client/gui"
)

func main() {
	env.LoadEnv()
	gui.Start()
}
