package client

import (
	"log"
	"net"

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

//RequestDownload - TODO
func (client *Client) RequestDownload(sessionID string) {
	client.checkConnected()
	lib.SendString(client.connection, "get", lib.COMMANDSIZE)
	lib.SendString(client.connection, sessionID, lib.COMMANDSIZE)
	lib.DownloadFile(client.connection)
}

//WaitAndSendFile - TODO
func (client *Client) WaitAndSendFile(filePath string) {
	client.checkConnected()
	response, err := lib.ReadString(client.connection, lib.COMMANDSIZE)
	lib.CheckFatalError(err)
	if response == "start" {
		lib.SendFile(client.connection, filePath)
	}
}
