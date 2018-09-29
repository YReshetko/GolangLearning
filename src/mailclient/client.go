package main

import (
	"fmt"
	"mailclient/config"
	"mailclient/save"
	"mailclient/service"
	"os"
)

//Run to hide windows console of the application
//go build -ldflags "-H windowsgui" -v -o client.exe client.go
func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	complete := make(chan int)
	var config config.Configuration

	if err := config.Load(); err != nil {
		fmt.Println("Can not run an application as config file was not loaded due to next error:", err)
		os.Exit(1)
	}
	fmt.Printf("Config:\n%+v\n", config)

	dbAccess := save.NewDBAccess(config.StorageConfiguration.DbHost, config.StorageConfiguration.DbPort, config.StorageConfiguration.DbName)
	dbAccess.StartSession()
	defer dbAccess.CloseSession()
	dao := save.NewDao(dbAccess.GetCollection(config.StorageConfiguration.CollectionName))

	emailService := service.NewEmailFetcher(config, dao)
	go service.Job(emailService, config.SchedulerConfiguration)
	go service.RunWebService(config.StorageConfiguration, emailService, dao)
	go service.StartAppInTray(complete)

	<-complete
}
