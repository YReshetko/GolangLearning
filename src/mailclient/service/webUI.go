package service

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"mailclient/config"
	"mailclient/domain"
	"mailclient/save"
	"mailclient/util"
	"net/http"
	"time"
)

var (
	dbAccess     save.DBAccess
	emailService EmailService
	dao          save.EmailDao
	fileStorage  string
)

func RunWebService(config config.StorageConfig, service EmailService, emailDao save.EmailDao) {
	emailService = service
	dbAccess = save.NewDBAccess(config.DbHost, config.DbPort, config.DbName)
	dao = emailDao
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
	latestRecords := dao.FindLatest(10)
	renderEmailData(w, latestRecords)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	date1 := r.FormValue("date1")
	date2 := r.FormValue("date2")
	callType := r.FormValue("callType")
	if date1 != "" || date2 != "" {
		from, to := util.GetDateRange(date1, date2)
		fmt.Printf("Serch for: date1:%v; date2:%v, type:%v\n", from, to, callType)
		records := dao.FindByDateRange(from, to)
		renderEmailData(w, records)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
func renderEmailData(w http.ResponseWriter, emailData []domain.EmailData) {
	t, _ := template.ParseFiles("web/index.html")
	t.Execute(w, emailData)
}

func startServer() error {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("web/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("web/img"))))
	http.Handle("/records/", http.StripPrefix("/records/", http.FileServer(http.Dir(fileStorage))))
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/search", searchHandler)
	return http.ListenAndServe(":8080", nil)
}
