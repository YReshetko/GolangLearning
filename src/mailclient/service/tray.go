package service

import (
	"io/ioutil"
	"log"
	"mailclient/util"
	"os"

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
	open := systray.AddMenuItem("Браузер", "Доступ к поиску")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")
	go func() {
		for {
			select {
			case <-open.ClickedCh:
				util.OpenWindow("http://localhost:8080")
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
	log.Println("Closing tray")
	os.Exit(0)
	close <- 1

}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Println("Error during loading tray icon:", err)
	}
	return b
}
