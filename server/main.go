package server

import (
	"log"
	"math/rand"
	"net"
	"strconv"

	"github.com/emil-nasso/share/conn"
	lib "github.com/emil-nasso/share/lib"
)

//Uploader - TODO
//TODO - Gotta remember to clear these out when they disconnect
type Uploader struct {
	sessionID       string
	connection      *conn.ShareConnection
	downloaders     chan *conn.ShareConnection
	downloadersHTTP chan *FileTransferRequest
}

func (uploader *Uploader) run() {
	for {
		//TODO - quit-channel?
		select {
		case downloader := <-uploader.downloaders:
			relayFileTransfer(uploader.connection, downloader)
		case fileTransferRequest := <-uploader.downloadersHTTP:
			relayHTTPTransfer(uploader.connection, fileTransferRequest.responseWriter)
			fileTransferRequest.done <- true
		}
	}
}

//Server - TODO
type Server struct {
	ip        string
	port      string
	uploaders []*Uploader
}

//New - Create a new server with standard configuration
func New() Server {
	return Server{ip: "localhost", port: "27001"}
}

//Start - Start the server and begin listening for new connections
func (server *Server) Start() error {
	go server.startHTTPServer()

	listener, err := net.Listen("tcp", server.ip+":"+server.port)
	lib.CheckFatalError(err)
	defer listener.Close()

	log.Println("File transfer server started. Waiting for connections.")
	for {
		connection, err := listener.Accept()
		lib.CheckError(err)

		log.Println("Client connected:", connection.RemoteAddr().String())
		go server.handleCommands(connection)
	}
}

func (server *Server) handleCommands(c net.Conn) {
	connection := conn.New(c)
	//	defer connection.Close()
	clientVersion, err := connection.ReadString(lib.COMMANDSIZE)
	if lib.CheckError(err) {
		return
	}
	connection.SendString(lib.PROTOCOLVERSION, lib.COMMANDSIZE)
	if clientVersion != lib.PROTOCOLVERSION {
		log.Println("Bad client version:", clientVersion, ". Server running version:", lib.PROTOCOLVERSION)
	}

commandloop:
	for {
		command, err := connection.ReadString(lib.COMMANDSIZE)
		if lib.CheckError(err) {
			break
		}
		switch command {
		case "upload":
			sessionID := generateSessionID()
			connection.SendString(sessionID, lib.COMMANDSIZE)
			uploader := Uploader{
				sessionID:       sessionID,
				connection:      connection,
				downloaders:     make(chan *conn.ShareConnection),
				downloadersHTTP: make(chan *FileTransferRequest),
			}
			server.uploaders = append(server.uploaders, &uploader)
			//Spawn a go routine here that checks if the connection is alive, until it's dead,
			// and remove it from uploaders.
			go uploader.run()
			break commandloop
		case "get":
			sessionID, err := connection.ReadString(lib.COMMANDSIZE)
			lib.CheckError(err)
			log.Println("Exchanging file for session", sessionID)
			uploader := server.findUploader(sessionID)
			uploader.downloaders <- connection
		}
	}
}

func (server *Server) findUploader(sessionID string) *Uploader {
	for _, uploader := range server.uploaders {
		if uploader.sessionID == sessionID {
			return uploader
		}
	}
	return nil
}

func generateSessionID() string {
	return strconv.Itoa(rand.Intn(899999999) + 100000000)
}
