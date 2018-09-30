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

const (
	urlRoot    = "/"
	urlSearch  = "/search"
	urlProcess = "/process"

	urlStaticCSS          = "/css/"
	urlStaticJS           = "/js/"
	urlStaticImage        = "/img/"
	urlStaticLocalStorage = "/records/"

	pathStaticRoot  = "web"
	pathStaticCSS   = pathStaticRoot + "/css"
	pathStaticJS    = pathStaticRoot + "/js"
	pathStaticImage = pathStaticRoot + "/img"

	pathPageEmailViewer   = pathStaticRoot + "/index.html"
	pathPageError         = pathStaticRoot + "/error.html"
	pathTemplatesRoot     = pathStaticRoot + "/tmp"
	pathHeaderTemplate    = pathTemplatesRoot + "/header.html"
	pathSearchTemplate    = pathTemplatesRoot + "/search.html"
	pathEmailViewTemplate = pathTemplatesRoot + "/emailoutput.html"

	templateIndex = "index"
	templateError = "error"
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
		http.Redirect(w, r, urlRoot, http.StatusFound)
	}
}
func processHandler(w http.ResponseWriter, r *http.Request) {
	if err := emailService.Process(); err != nil {
		renderErrorPage(w, err)
	} else {
		http.Redirect(w, r, urlRoot, http.StatusFound)
	}
}

func renderErrorPage(w http.ResponseWriter, err error) {
	t, _ := template.ParseFiles(pathPageError, pathHeaderTemplate)
	t.ExecuteTemplate(w, templateError, err)
}

func renderEmailData(w http.ResponseWriter, emailData []domain.EmailData) {
	t, _ := template.ParseFiles(pathPageEmailViewer,
		pathHeaderTemplate, pathSearchTemplate, pathEmailViewTemplate)
	t.ExecuteTemplate(w, templateIndex, emailData)
}

func startServer() error {
	http.Handle(urlStaticCSS, http.StripPrefix(urlStaticCSS, http.FileServer(http.Dir(pathStaticCSS))))
	http.Handle(urlStaticJS, http.StripPrefix(urlStaticJS, http.FileServer(http.Dir(pathStaticJS))))
	http.Handle(urlStaticImage, http.StripPrefix(urlStaticImage, http.FileServer(http.Dir(pathStaticImage))))
	http.Handle(urlStaticLocalStorage, http.StripPrefix(urlStaticLocalStorage, http.FileServer(http.Dir(fileStorage))))
	http.HandleFunc(urlRoot, welcomeHandler)
	http.HandleFunc(urlSearch, searchHandler)
	http.HandleFunc(urlProcess, processHandler)
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
