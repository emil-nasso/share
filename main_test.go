package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/emil-nasso/share/client"
	"github.com/emil-nasso/share/server"
)

func TestFullStack(t *testing.T) {
	//Start the server
	server := server.New()
	go server.Start()
	//HACK: return a channel
	time.Sleep(time.Second)

	timestamp := string(time.Now().Unix())
	uploaderDirectory := "/tmp/share-test-data/" + timestamp + "/uploader/"
	uploaderFile := uploaderDirectory + "/test.txt"
	os.MkdirAll(uploaderDirectory, 0777)
	ioutil.WriteFile(uploaderFile, []byte(timestamp), 0777)

	uploader := client.New("127.0.0.1")
	sessionID := uploader.RequestUpload()
	defer uploader.Disconnect()
	go uploader.WaitAndSendFile(uploaderFile)
	//HACK: return a channel
	time.Sleep(time.Second)

	downloader := client.New("127.0.0.1")
	defer downloader.Disconnect()
	downloader.RequestDownload(sessionID)
}
