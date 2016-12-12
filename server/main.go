package server

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"

	lib "github.com/emil-nasso/share/lib"
)

//Uploader - TODO
//TODO - Gotta remember to clear these out when they disconnect
type Uploader struct {
	sessionID string
	uploader  net.Conn
}

//Server - TODO
type Server struct {
	ip        string
	port      string
	uploaders []Uploader
}

//New - Create a new server with standard configuration
func New() Server {
	return Server{ip: "localhost", port: "27001"}
}

//Start - Start the server and begin listening for new connections
func (server *Server) Start() error {
	listener, err := net.Listen("tcp", server.ip+":"+server.port)
	lib.CheckFatalError(err)
	defer listener.Close()

	log.Println("Server started. Waiting for connections...")
	for {
		connection, err := listener.Accept()
		lib.CheckError(err)

		log.Println("Client connected from: ", connection.RemoteAddr().String())
		go server.handleCommands(connection)
	}
}

func (server *Server) handleCommands(connection net.Conn) {
	//	defer connection.Close()
	clientVersion, err := lib.ReadString(connection, lib.COMMANDSIZE)
	if lib.CheckError(err) {
		return
	}
	lib.SendString(connection, lib.PROTOCOLVERSION, lib.COMMANDSIZE)
	if clientVersion != lib.PROTOCOLVERSION {
		log.Println("Bad client version:", clientVersion, ". Server running version: ", lib.PROTOCOLVERSION)
	}

commandloop:
	for {
		command, err := lib.ReadString(connection, lib.COMMANDSIZE)
		if lib.CheckError(err) {
			break
		}
		switch command {
		case "hello":
			lib.SendString(connection, "zup", lib.COMMANDSIZE)
		case "sendfile":
			server.downloadFile(connection)
		case "upload":
			sessionID := generateSessionID()
			lib.SendString(connection, sessionID, lib.COMMANDSIZE)
			server.uploaders = append(server.uploaders, Uploader{sessionID: sessionID, uploader: connection})
			//Spawn a go routine here that checks if the connection is alive, until it's dead,
			// and remove it from uploaders.
			break commandloop
		}
	}
}

func (server *Server) downloadFile(connection net.Conn) {
	fileName, err := lib.ReadString(connection, lib.COMMANDSIZE)
	lib.CheckError(err)
	fileSizeData, err := lib.ReadString(connection, lib.COMMANDSIZE)
	lib.CheckError(err)
	fileSize, _ := strconv.ParseInt(fileSizeData, 10, 64)

	newFile, err := os.Create(fileName)
	lib.CheckError(err)

	defer newFile.Close()

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < lib.BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+lib.BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, lib.BUFFERSIZE)
		receivedBytes += lib.BUFFERSIZE
	}
	log.Println("Received file completely!")
}

func generateSessionID() string {
	return "this-is-not-very-random"
}
