package service

import (
	"html/template"
	"io/ioutil"
	"mailclient/config"
	"mailclient/logger"
	"mailclient/save"
	"mailclient/util"
	"net/http"
	"time"
)

const (
	urlRoot        = "/"
	urlSearch      = "/search"
	urlProcess     = "/process"
	urlDiagnostic  = "/diagnostic"
	urlFixDaoIssue = "/fix/dao"

	urlStaticCSS          = "/css/"
	urlStaticJS           = "/js/"
	urlStaticImage        = "/img/"
	urlStaticLocalStorage = "/records/"

	pathStaticRoot  = "web"
	pathStaticCSS   = pathStaticRoot + "/css"
	pathStaticJS    = pathStaticRoot + "/js"
	pathStaticImage = pathStaticRoot + "/img"

	pathPageEmailViewer    = pathStaticRoot + "/index.html"
	pathPageDiagnostic     = pathStaticRoot + "/diagnostic.html"
	pathTemplatesRoot      = pathStaticRoot + "/tmp"
	pathHeaderTemplate     = pathTemplatesRoot + "/header.html"
	pathSearchTemplate     = pathTemplatesRoot + "/search.html"
	pathEmailViewTemplate  = pathTemplatesRoot + "/emailoutput.html"
	pathDiagnosticTemplate = pathTemplatesRoot + "/diagnosticStatus.html"

	templateIndex      = "index"
	templateDiagnostic = "diagnostic"
)

var (
	emailService      EmailService
	dao               save.EmailDao
	fileStorage       string
	diagnosticService DiagnosticService
)

type DiagnosticStatus struct {
	ImapSt  ImapStatus
	ImapErr error
	DaoSt   DaoStatus
	DaoErr  error
	DiscSt  StorageStatus
	DiscErr error
}
type viewModel struct {
	Error error
	Data  interface{}
}

/*
RunWebService - run web service
*/
func RunWebService(config config.StorageConfig, service EmailService, emailDao save.EmailDao, diagnostic DiagnosticService) {
	emailService = service
	dao = emailDao
	fileStorage = config.LocalStorageBasePath
	diagnosticService = diagnostic
	for {
		logger.Info("Starting web service")
		err := startServer()
		if err != nil {
			logger.Error("Web service is crashed with error:", err)
			logger.Error("Restarting web service in 1 minute")
			time.Sleep(time.Minute)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	latestRecords, err := dao.FindLatest(10)
	model := viewModel{}
	if err != nil {
		model.Error = err
	} else {
		model.Data = latestRecords
	}
	renderEmailData(w, model)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	date1 := r.FormValue("date1")
	date2 := r.FormValue("date2")
	if date1 != "" || date2 != "" {
		from, to := util.GetDateRange(date1, date2)
		logger.Debug("Serch for: date1:%v; date2:%v\n", from, to)
		records, err := dao.FindByDateRange(from, to)
		model := viewModel{}
		if err != nil {
			model.Error = err
		} else {
			model.Data = records
		}
		renderEmailData(w, model)
	} else {
		http.Redirect(w, r, urlRoot, http.StatusFound)
	}
}
func processHandler(w http.ResponseWriter, r *http.Request) {
	if err := emailService.Process(); err != nil {
		model := viewModel{}
		model.Error = err
		renderEmailData(w, model)
	} else {
		http.Redirect(w, r, urlRoot, http.StatusFound)
	}
}

func diagnosticHandler(w http.ResponseWriter, r *http.Request) {
	status := &DiagnosticStatus{}
	status.ImapSt, status.ImapErr = diagnosticService.CheckImap()
	status.DaoSt, status.DaoErr = diagnosticService.CheckDao()
	status.DiscSt, status.DiscErr = diagnosticService.CheckLocalStorage()
	t, err := template.ParseFiles(pathPageDiagnostic, pathHeaderTemplate, pathDiagnosticTemplate)
	logger.Debug("Error rendering diagnostic page:", err)
	model := viewModel{err, status}
	t.ExecuteTemplate(w, templateDiagnostic, model)
}

func fixDaoHandler(w http.ResponseWriter, r *http.Request) {
	diagnosticService.FixDao()
	http.Redirect(w, r, urlDiagnostic, http.StatusFound)
}

func renderEmailData(w http.ResponseWriter, model viewModel) {
	t, _ := template.ParseFiles(pathPageEmailViewer,
		pathHeaderTemplate, pathSearchTemplate, pathEmailViewTemplate)
	t.ExecuteTemplate(w, templateIndex, model)
}

func startServer() error {
	http.Handle(urlStaticCSS, http.StripPrefix(urlStaticCSS, http.FileServer(http.Dir(pathStaticCSS))))
	http.Handle(urlStaticJS, http.StripPrefix(urlStaticJS, http.FileServer(http.Dir(pathStaticJS))))
	http.Handle(urlStaticImage, http.StripPrefix(urlStaticImage, http.FileServer(http.Dir(pathStaticImage))))
	http.Handle(urlStaticLocalStorage, http.StripPrefix(urlStaticLocalStorage, http.FileServer(http.Dir(fileStorage))))
	http.HandleFunc(urlRoot, welcomeHandler)
	http.HandleFunc(urlSearch, searchHandler)
	http.HandleFunc(urlProcess, processHandler)
	http.HandleFunc(urlDiagnostic, diagnosticHandler)
	http.HandleFunc(urlFixDaoIssue, fixDaoHandler)
	return http.ListenAndServe(":8080", nil)
}

func loadFile(fileName string) []byte {
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error("Error during file loading %s, error: %v", fileName, err)
		return nil
	}
	return body
}
