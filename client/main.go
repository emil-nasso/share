package client

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/emil-nasso/share/conn"
	"github.com/emil-nasso/share/ftpb"
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
		fmt.Println("Starting upload")
		c.uploadFile(filePath)
	}
	fmt.Println("Done with upload")
}

func (c *Client) downloadFile() string {
	fileName, fileSizeData := c.serverConnection.GetFileNameAndSize()
	fileSize, _ := strconv.ParseInt(fileSizeData, 10, 64)

	bar := ftpb.New(int64(fileSize))
	bar.Start()

	newFile, err := os.Create(fileName)
	lib.CheckError(err)

	defer newFile.Close()

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) <= lib.BUFFERSIZE {
			bytesLeft := (fileSize - receivedBytes)
			io.CopyN(newFile, c.serverConnection, bytesLeft)
			c.serverConnection.Read(make([]byte, lib.BUFFERSIZE-bytesLeft))
			bar.Done()
			fmt.Println("Download completed")
			break
		}
		io.CopyN(newFile, c.serverConnection, lib.BUFFERSIZE)
		receivedBytes += lib.BUFFERSIZE
		bar.Set(int64(receivedBytes))
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
	bar := ftpb.New(int64(fileSizeData))
	bar.Start()

	c.serverConnection.SendFileNameAndSize(fileName, fileSize)

	sendBuffer := make([]byte, lib.BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		bar.Increase(lib.BUFFERSIZE)
		c.serverConnection.Write(sendBuffer)
	}
	bar.Done()
	fmt.Println("Transfer complete")
	return
}
