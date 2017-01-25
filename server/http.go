package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
)

//FileTransferRequest - TODO
type FileTransferRequest struct {
	responseWriter http.ResponseWriter
	done           chan bool
}

func (server *Server) startHTTPServer() {
	r := mux.NewRouter()
	r.HandleFunc("/status", (&templateHandler{filename: "status.html"}).serveHTTP)
	r.HandleFunc("/get/{session-id}", server.getFileHandler)
	http.Handle("/", r)
	log.Println("Starting http server")
	go http.ListenAndServe(":27002", nil)
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, struct {
		Title string
	}{Title: "This is a statuspage"})
}

func (server *Server) getFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["session-id"]

	uploader := server.findUploader(sessionID)
	if uploader == nil {
		http.NotFound(w, r)
		return
	}
	done := make(chan bool)
	transferRequest := FileTransferRequest{
		responseWriter: w,
		done:           done,
	}

	uploader.downloadersHTTP <- &transferRequest
	<-done
}
