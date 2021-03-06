package service

import (
	"io/ioutil"
	"mailclient/logger"
	"mailclient/util"

	"github.com/getlantern/systray"
)

var (
	close chan int
)

/*
StartAppInTray - starts application in tray
*/
func StartAppInTray(blockingChan chan int) {
	close = blockingChan
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("icon.ico"))
	systray.SetTitle("E-MFetcher")
	systray.SetTooltip("E-MFetcher - сбор писем")
	about := systray.AddMenuItem("О программе", "Информация о программе")
	logMenu := systray.AddMenuItem("Лог", "Открыть файл с логом")
	open := systray.AddMenuItem("Браузер", "Доступ к поиску")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")
	go func() {
		for {
			select {
			case <-open.ClickedCh:
				util.OpenWindow("http://localhost:8080")
			case <-logMenu.ClickedCh:
				util.OpenWindow("/log/mfetch.log")
			case <-about.ClickedCh:
				util.OpenWindow("about.txt")
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	logger.Info("Closing tray")
	close <- 1
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		logger.Error("Error during loading tray icon:", err)
	}
	return b
}
