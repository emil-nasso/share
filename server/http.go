package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//FileTransferRequest - TODO
type FileTransferRequest struct {
	responseWriter http.ResponseWriter
	done           chan bool
}

func (server *Server) startHTTPServer() {
	r := mux.NewRouter()
	r.HandleFunc("/status", server.statusHandler)
	r.HandleFunc("/get/{session-id}", server.getFileHandler)
	http.Handle("/", r)
	log.Println("Starting http server")
	go http.ListenAndServe(":27002", nil)
}

func (server *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Statuspage!</h1>")
	fmt.Fprint(w, "<h3>Sessions</h3>")
	fmt.Fprint(w, "<ul>")
	for _, uploader := range server.uploaders {
		fmt.Fprintf(w, "<li>%v</li>", uploader.sessionID)
	}
	fmt.Fprint(w, "</ul>")
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
