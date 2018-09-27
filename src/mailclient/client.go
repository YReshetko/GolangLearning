package main

import (
	"fmt"
	"mailclient/config"
	"mailclient/service"
	"os"
)

func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	var config config.Configuration
	if err := config.Load(); err != nil {
		fmt.Println("Can not run an application as config file was not loaded due to next error:", err)
		os.Exit(1)
	}
	fmt.Printf("Config:\n%+v\n", config)
	emailService := service.NewEmailFetcher(config)
	go service.Job(emailService, config.SchedulerConfiguration)

	//time.Sleep(2 * time.Minute)
	//for {} // - hangup of application at imapClient.go -> bufferCompleted <- fetchManager.FetchFunction()(seqset, fetchManager.FetchItems(), messages) on second try
	select {}
	//complete := make(chan error)
	//<-complete
}
