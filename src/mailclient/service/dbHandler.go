package service

import (
	"log"
	"mailclient/config"
	"mailclient/util"
)

var mongoExe = "mongod.exe"

type mongoHandler struct {
	mongoAppPath string
	mongoDbPath  string
}
type DbHandler interface {
	Start() error
	Restart() error
	Stop() error
}

func NewDbHandler(config config.StorageConfig) DbHandler {
	return &mongoHandler{
		config.MongoAppPath + mongoExe,
		config.MongoDbPath,
	}
}

func (handler *mongoHandler) Start() error {
	pid, err := util.FindPIDByName(mongoExe)
	if err != nil {
		log.Printf("Error searching %s pid process: %v\n", mongoExe, err)
	}
	if pid > 0 {
		log.Printf("Process %s with pid %v is already exist\n", mongoExe, pid)
	} else {
		util.RunWinProgramm(handler.mongoAppPath, []string{"--dbpath", handler.mongoDbPath})
	}
	return err
}
func (handler *mongoHandler) Stop() error {
	pid, err := util.FindPIDByName(mongoExe)
	if err != nil {
		log.Printf("Error searching %s pid process: %v\n", mongoExe, err)
	}
	if pid > 0 {
		util.KillPid(pid)
	} else {
		log.Println("Can't kill mongo PID because it doesn't exist")
	}
	return err
}

func (handler *mongoHandler) Restart() error {
	handler.Stop()
	return handler.Start()
}
