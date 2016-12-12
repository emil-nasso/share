package client

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"

	lib "github.com/emil-nasso/share/lib"
)

//Client - TODO
type Client struct {
	ip         string
	port       string
	connection net.Conn
}

//New - TODO
func New() Client {
	return Client{ip: "localhost", port: "27001", connection: nil}
}

//Connect - TODO
func (client *Client) Connect() {
	connection, err := net.Dial("tcp", client.ip+":"+client.port)
	lib.CheckFatalError(err)
	client.connection = connection
	client.negotiateVersion()
}

//Disconnect - TODO
func (client *Client) Disconnect() {
	client.connection.Close()
}

func (client *Client) negotiateVersion() {
	client.checkConnected()
	lib.SendString(client.connection, lib.PROTOCOLVERSION, lib.COMMANDSIZE)
	serverVersion, err := lib.ReadString(client.connection, lib.COMMANDSIZE)
	if err != nil {
		log.Fatalln(err)
	}
	if serverVersion != lib.PROTOCOLVERSION {
		log.Fatalln("Client/server version mismatch. Server: ", serverVersion, "Client: ", lib.PROTOCOLVERSION)
	}
}

func (client *Client) checkConnected() {
	if client.connection == nil {
		panic("Not connected")
	}
}

//RequestUpload - TODO
func (client *Client) RequestUpload() string {
	client.checkConnected()
	lib.SendString(client.connection, "upload", lib.COMMANDSIZE)
	sessionID, err := lib.ReadString(client.connection, lib.COMMANDSIZE)
	lib.CheckFatalError(err)
	return sessionID
}

//WaitAndSendFile - TODO
func (client *Client) WaitAndSendFile(filePath string) {
	client.checkConnected()
	response, err := lib.ReadString(client.connection, lib.COMMANDSIZE)
	lib.CheckFatalError(err)
	if response == "start" {
		client.sendFile(filePath)
	}
}

func (client *Client) sendFile(filePath string) {
	client.checkConnected()
	file, err := os.Open(filePath)
	lib.CheckFatalError(err)
	fileInfo, err := file.Stat()
	lib.CheckFatalError(err)

	fileSize := strconv.FormatInt(fileInfo.Size(), 10)
	fileName := fileInfo.Name()

	lib.SendString(client.connection, "sendfile", lib.COMMANDSIZE)
	lib.SendString(client.connection, fileName, lib.COMMANDSIZE)
	lib.SendString(client.connection, fileSize, lib.COMMANDSIZE)

	log.Println("Transfering file")
	sendBuffer := make([]byte, lib.BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		client.connection.Write(sendBuffer)
	}
	log.Println("Transfer complete")
	return
}
