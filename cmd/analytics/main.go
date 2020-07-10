package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, cleanup, err := wireUp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	handleInterruptSignal(func(c chan os.Signal) {
		<-c
		fmt.Println("shutting down...")
	})

	app.jobRunner.Start()

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("http server terminated: %s", err.Error()))
		}
	}()

	go func() {
		app.workerPool.Run()
	}()

	if err := app.messageConsumer.ConsumeMessages(); err != nil {
		panic(fmt.Sprintf("message consumer terminated: %s", err.Error()))
	}

	fmt.Printf("shut down completed!")
}

func handleInterruptSignal(callback func(c chan os.Signal)) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go callback(c)
}
