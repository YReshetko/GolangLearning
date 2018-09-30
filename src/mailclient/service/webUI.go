package service

import (
	"html/template"
	"io/ioutil"
	"log"
	"mailclient/config"
	"mailclient/domain"
	"mailclient/save"
	"mailclient/util"
	"net/http"
	"time"
)

var (
	emailService EmailService
	dao          save.EmailDao
	fileStorage  string
)

/*
RunWebService - run web service
*/
func RunWebService(config config.StorageConfig, service EmailService, emailDao save.EmailDao) {
	emailService = service
	dao = emailDao
	fileStorage = config.LocalStorageBasePath
	for {
		log.Println("Starting web service")
		err := startServer()
		if err != nil {
			log.Println("Web service is crashed with error:", err)
			log.Println("Restarting web service in 1 minute")
			time.Sleep(time.Minute)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	latestRecords := dao.FindLatest(10)
	renderEmailData(w, latestRecords)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	date1 := r.FormValue("date1")
	date2 := r.FormValue("date2")
	if date1 != "" || date2 != "" {
		from, to := util.GetDateRange(date1, date2)
		log.Printf("Serch for: date1:%v; date2:%v\n", from, to)
		records := dao.FindByDateRange(from, to)
		renderEmailData(w, records)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func renderEmailData(w http.ResponseWriter, emailData []domain.EmailData) {
	t, _ := template.ParseFiles("web/index.html",
		"web/tmp/header.html", "web/tmp/search.html", "web/tmp/emailoutput.html")
	t.ExecuteTemplate(w, "index", emailData)
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

func loadFile(fileName string) []byte {
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error during file loading %s, error: %v", fileName, err)
		return nil
	}
	return body
}
