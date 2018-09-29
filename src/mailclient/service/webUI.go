package service

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"mailclient/config"
	"mailclient/save"
	"net/http"
	"time"
)

var (
	dbAccess       save.DBAccess
	emailService   EmailService
	collectionName string
	fileStorage    string
)

func RunWebService(config config.StorageConfig, service EmailService) {
	emailService = service
	dbAccess = save.NewDBAccess(config.DbHost, config.DbPort, config.DbName)
	collectionName = config.CollectionName
	fileStorage = config.LocalStorageBasePath
	for {
		err := startServer()
		if err != nil {
			fmt.Println("Web service is crashed with error:", err)
			fmt.Println("Restarting web service in 1 minute")
			time.Sleep(time.Minute)
		}
	}
}

func loadFile(title string) ([]byte, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	/*indexPage, err := loadFile("web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
	dbAccess.StartSession()
	defer dbAccess.CloseSession()
	collection := dbAccess.GetCollection(collectionName)
	dao := save.NewDao(collection)
	latestRecords := dao.FindLatest(10)

	t, _ := template.ParseFiles("web/index.html")
	t.Execute(w, latestRecords)
	//w.Write(indexPage)
}

func startServer() error {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("web/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("web/img"))))
	http.HandleFunc("/", welcomeHandler)
	return http.ListenAndServe(":8080", nil)
}
