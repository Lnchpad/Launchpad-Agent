package main

import (
	"cjavellana.me/launchpad/agent/app"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	sigChannel := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	go func() {
		sig := <-sigChannel
		fmt.Println(sig)
		done <- true
	}()

	agent := app.NewAgent()
	agent.Start()
	defer agent.Terminate()

	fmt.Println("Press Ctrl + C to Terminate")
	<-done
	fmt.Println("Exiting")
}
