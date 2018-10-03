package main

import (
	"log"
	"mailclient/config"
	"mailclient/save"
	"mailclient/service"
	"os"
	"time"
)

//Run to hide windows console of the application
//go build -ldflags "-H windowsgui" -v -o EMFetcher.exe client.go
func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	close := initLogging()
	defer close()
	complete := make(chan int)
	var config config.Configuration
	if err := config.Load(); err != nil {
		log.Println("Can not run an application as config file was not loaded due to next error:", err)
		os.Exit(1)
	}
	dbHandler := service.NewDbHandler(config.StorageConfiguration)
	err := dbHandler.Start()
	if err != nil {
		log.Println("DB was not started at application start due to ", err)
	} else {
		defer dbHandler.Stop()
		time.Sleep(2 * time.Second)
	}
	log.Printf("Start init DB access: %+v\n", config.StorageConfiguration)
	dbAccess := save.NewDBAccess(config.StorageConfiguration.DbHost, config.StorageConfiguration.DbPort, config.StorageConfiguration.DbName)
	dbAccess.StartSession()
	defer dbAccess.CloseSession()
	dao := save.NewDao(dbAccess.GetCollection(config.StorageConfiguration.CollectionName))

	log.Printf("Init fetch email service: %+v\n", config.EmailStructure)
	emailService := service.NewEmailFetcher(config, dao)

	diagnosticService := service.NewDiagnosticService(emailService, dbHandler, dao, dbAccess, config)
	log.Println("Starting services")
	go service.Job(emailService, config.SchedulerConfiguration)
	go service.RunWebService(config.StorageConfiguration, emailService, dao, diagnosticService)
	go service.StartAppInTray(complete)

	log.Println("Starting process as a first fetch at application start")
	//emailService.Process()

	<-complete
}

func initLogging() func() {
	f, err := os.OpenFile("log/mfetch.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	log.Println("Start application")
	return func() {
		log.Println("End application")
		f.Close()
	}
}
