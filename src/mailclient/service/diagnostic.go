package service

import (
	"fmt"
	"mailclient/config"
	"mailclient/logger"
	"mailclient/save"
	"mailclient/util"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type diagnostic struct {
	imapService EmailService
	dbHandler   DbHandler
	dao         save.EmailDao
	dbAccess    save.DBAccess
	config      config.Configuration
}
type GeneralStatus struct {
	Message string
	Status  bool
}

type ImapStatus struct {
	GeneralStatus
	Host string
	Port int
}

type DaoStatus struct {
	GeneralStatus
	Host string
	Port string
}
type StorageStatus struct {
	GeneralStatus
	TotalSpace string
	UsedSpace  string
	FreeSpace  string
}

type DiagnosticService interface {
	CheckImap() (ImapStatus, error)
	CheckDao() (DaoStatus, error)
	CheckLocalStorage() (StorageStatus, error)

	FixDao()
}

func NewDiagnosticService(imap EmailService, dbHandler DbHandler, dao save.EmailDao, access save.DBAccess, appConfig config.Configuration) DiagnosticService {
	return &diagnostic{imap, dbHandler, dao, access, appConfig}
}

func (diag *diagnostic) CheckImap() (ImapStatus, error) {
	status := ImapStatus{
		GeneralStatus{},
		diag.config.HostConfiguration.ImapHost,
		diag.config.HostConfiguration.ImapPort,
	}
	err := diag.imapService.PrintMailboxes()
	if err != nil {
		status.Status = false
		status.Message = "Проверьте интернет соединение и доступность почты через web browser"
	} else {
		status.Status = true
		status.Message = "Почта доступна"
	}
	return status, err
}

func (diag *diagnostic) CheckDao() (DaoStatus, error) {
	status := DaoStatus{
		GeneralStatus{},
		diag.config.StorageConfiguration.DbHost,
		diag.config.StorageConfiguration.DbPort,
	}
	_, err := dao.FindLatest(1)
	if err != nil {
		status.Status = false
		status.Message = "База данных недоступна, попробуйте перезапустить БД, нажать кнопку \"Исправить\" и провести диагностику снова. Если не поможет, то перезапустить приложение"
	} else {
		status.Status = true
		status.Message = "База данных доступна"
	}
	return status, nil
}

func (diag *diagnostic) CheckLocalStorage() (StorageStatus, error) {
	status := StorageStatus{}
	disk := util.DiskUsage(diag.config.StorageConfiguration.LocalStorageBasePath)
	freeSpace := float64(disk.Free) / float64(GB)
	totalSpace := float64(disk.All) / float64(GB)
	status.TotalSpace = fmt.Sprintf("%.2f GB\n", totalSpace)
	status.UsedSpace = fmt.Sprintf("%.2f GB\n", float64(disk.Used)/float64(GB))
	status.FreeSpace = fmt.Sprintf("%.2f GB\n", freeSpace)

	if freeSpace < 2 || (freeSpace*100/totalSpace) < 5 {
		status.Status = false
		status.Message = "Осталось мало свободного места на жестком диске, где расположено хранилище (<2Gb или <5%). Рекомендация: очистить диск"
	} else {
		status.Status = true
		status.Message = "Памяти под хранилище достаточно для продолжения работы"
	}
	return status, nil
}

func (diag *diagnostic) FixDao() {
	diag.dbAccess.CloseSession()
	err := diag.dbHandler.Restart()
	if err != nil {
		logger.Error("Error during restarting DB:", err)
	}
	time.Sleep(2 * time.Second)
	diag.dbAccess.StartSession()
	collection := diag.dbAccess.GetCollection(diag.config.StorageConfiguration.CollectionName)
	diag.dao.UpdateCollection(collection)
}
