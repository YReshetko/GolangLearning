package service

import (
	"log"
	"mailclient/config"
	"time"

	"github.com/jasonlvhit/gocron"
)

var executeVar EmailService

/*
Job - starts scheduler to fetch emails
*/
func Job(emailService EmailService, config config.SchedulerConfig) {
	job := gocron.Every(config.Every)
	executeVar = emailService
	switch config.Term {
	case "Second":
		if config.Every > 1 {
			job = job.Seconds()
		} else {
			job = job.Second()
		}
	case "Minute":
		if config.Every > 1 {
			job = job.Minutes()
		} else {
			job = job.Minute()
		}
	case "Hour":
		if config.Every > 1 {
			job = job.Hours()
		} else {
			job = job.Hour()
		}
	case "Day":
		if config.Every > 1 {
			job = job.Days()
		} else {
			job = job.Day()
		}
	}

	if config.At != "" && config.Term == "Day" {
		job = job.At(config.At)
	}

	job.Do(run)
	<-gocron.Start()
}

func run() {
	go func() {
		log.Println("Running fetch processing by scheduler at:", time.Now())
		if err := executeVar.Process(); err != nil {
			log.Println("Error during running fetch process by scheduler:", err)
		}
	}()
}
