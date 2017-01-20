package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emil-nasso/share/lib"
	"github.com/gorilla/mux"
)

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

	uploader := server.findUploaderConnection(sessionID)
	if uploader == nil {
		http.NotFound(w, r)
		return
	}
	lib.RelayHTTPTransfer(uploader, w)
}
