package service

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
)

var (
	close chan int
)

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
				openBrowser("http://localhost:8080")
			case <-about.ClickedCh:

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func onExit() {
	close <- 1
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
