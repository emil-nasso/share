package client

import (
	"io"
	"os"
	"strconv"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/emil-nasso/share/conn"
	lib "github.com/emil-nasso/share/lib"
)

//Client - TODO
type Client struct {
	ServerHostname   string
	port             string
	serverConnection *conn.ShareConnection
}

//New - TODO
func New(url string) Client {
	port := "27001"
	serverConnection := conn.Connect(url, port)
	return Client{ServerHostname: url, port: port, serverConnection: serverConnection}
}

//Disconnect - TODO
func (c *Client) Disconnect() {
	c.serverConnection.Disconnect()
}

//RequestUpload - TODO
func (c *Client) RequestUpload() string {
	c.serverConnection.CheckConnected()
	c.serverConnection.SendString("upload", lib.COMMANDSIZE)
	sessionID, err := c.serverConnection.ReadString(lib.COMMANDSIZE)
	lib.CheckFatalError(err)
	return sessionID
}

//RequestDownload - TODO
func (c *Client) RequestDownload(sessionID string) {
	c.serverConnection.CheckConnected()
	c.serverConnection.SendString("get", lib.COMMANDSIZE)
	c.serverConnection.SendString(sessionID, lib.COMMANDSIZE)
	c.downloadFile()
}

//WaitAndSendFile - TODO
func (c *Client) WaitAndSendFile(filePath string) {
	c.serverConnection.CheckConnected()
	response, err := c.serverConnection.ReadString(lib.COMMANDSIZE)
	lib.CheckFatalError(err)
	if response == "start" {
		c.uploadFile(filePath)
	}
}

func (c *Client) downloadFile() string {
	fileName, fileSizeData := c.serverConnection.GetFileNameAndSize()
	fileSize, _ := strconv.ParseInt(fileSizeData, 10, 64)

	bar := pb.StartNew(int(fileSize))

	newFile, err := os.Create(fileName)
	lib.CheckError(err)

	defer newFile.Close()

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < lib.BUFFERSIZE {
			io.CopyN(newFile, c.serverConnection, (fileSize - receivedBytes))
			c.serverConnection.Read(make([]byte, (receivedBytes+lib.BUFFERSIZE)-fileSize))
			bar.Set(int(fileSize))
			bar.FinishPrint("Download completed")
			break
		}
		io.CopyN(newFile, c.serverConnection, lib.BUFFERSIZE)
		receivedBytes += lib.BUFFERSIZE
		bar.Set(int(receivedBytes))
	}
	return fileName
}

func (c *Client) uploadFile(filePath string) {
	file, err := os.Open(filePath)
	lib.CheckFatalError(err)
	fileInfo, err := file.Stat()
	lib.CheckFatalError(err)

	fileSizeData := fileInfo.Size()
	fileSize := strconv.FormatInt(fileSizeData, 10)
	fileName := fileInfo.Name()
	bar := pb.StartNew(int(fileSizeData))

	c.serverConnection.SendFileNameAndSize(fileName, fileSize)

	sendBuffer := make([]byte, lib.BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			bar.Set(int(fileSizeData))
			bar.FinishPrint("Transfer complete")
			break
		}
		bar.Add(lib.BUFFERSIZE)
		c.serverConnection.Write(sendBuffer)
	}
	return
}
