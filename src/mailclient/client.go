package main

import (
	"mailclient/config"
	"mailclient/logger"
	"mailclient/save"
	"mailclient/service"
	"os"
	"time"
)

//Run to hide windows console of the application
//go build -ldflags "-H windowsgui" -v -o EMFetcher.exe client.go
func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	close := logger.Init()
	defer close()
	logger.SetLogLevel(logger.DEBUG)
	complete := make(chan int)
	var config config.Configuration
	if err := config.Load(); err != nil {
		logger.Error("Can not run an application as config file was not loaded due to next error:", err)
		os.Exit(1)
	}
	dbHandler := service.NewDbHandler(config.StorageConfiguration)
	err := dbHandler.Start()
	if err != nil {
		logger.Error("DB was not started at application start due to ", err)
	} else {
		defer dbHandler.Stop()
		time.Sleep(2 * time.Second)
	}
	logger.Info("Start init DB access: %+v\n", config.StorageConfiguration)
	dbAccess := save.NewDBAccess(config.StorageConfiguration.DbHost, config.StorageConfiguration.DbPort, config.StorageConfiguration.DbName)
	dbAccess.StartSession()
	defer dbAccess.CloseSession()
	dao := save.NewDao(dbAccess.GetCollection(config.StorageConfiguration.CollectionName))

	logger.Info("Init fetch email service: %+v\n", config.EmailStructure)
	emailService := service.NewEmailFetcher(config, dao)

	diagnosticService := service.NewDiagnosticService(emailService, dbHandler, dao, dbAccess, config)
	logger.Info("Starting services")
	go service.Job(emailService, config.SchedulerConfiguration)
	go service.RunWebService(config.StorageConfiguration, emailService, dao, diagnosticService)
	go service.StartAppInTray(complete)

	logger.Info("Starting process as a first fetch at application start")
	//emailService.Process()

	<-complete
}
