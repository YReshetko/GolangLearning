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
	go service.RunWebService(config.StorageConfiguration, emailService)
	complete := make(chan error)
	<-complete
}
