package lib

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

//BUFFERSIZE - The size of the buffer for sending files. This amount of bytes will be
//sent per chunk.
const BUFFERSIZE = 1024

//COMMANDSIZE - The default size of the chucks for sending command and for the responses
const COMMANDSIZE = 64

//PROTOCOLVERSION - the current client/server communication protocol version
const PROTOCOLVERSION = "v1"

//DebugEnabled - Is debug enabled?
var DebugEnabled bool

//TODO - there has to be a better way
func fillString(message string, toLength int) string {
	for {
		stringLength := len(message)
		if stringLength < toLength {
			message = message + ":"
			continue
		}
		break
	}
	return message
}

//TODO: http://golangtutorials.blogspot.se/2011/06/inheritance-and-subclassing-in-go-or.html
// lägg ut send och read i en struct, en communicator, använd "arv" för att både servern och klienten ska ha tillgång till den

//ReadString - Reads a string from the connection
func ReadString(conn net.Conn, length int) (string, error) {
	buffer := make([]byte, length)
	_, err := conn.Read(buffer)
	if err == io.EOF {
		return "", errors.New("Connection closed")
	} else if err != nil {
		return "", err
	}
	response := strings.Trim(string(buffer), ":")
	Debug("Read string:" + response)
	return response, nil
}

//SendString - sends a string to the connection
func SendString(conn net.Conn, str string, length int) {
	paddedString := fillString(str, length)
	Debug("Sent string: " + str)
	conn.Write([]byte(paddedString))
}

//CheckFatalError - checks if the error is non-nil and logs a fatal error (exiting)
func CheckFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

//CheckError - Checks if the error is non-nil, prints the error to console and returns
//true if there was an error and false if there wasn't
func CheckError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

//RelayFileTransfer - TODO
func RelayFileTransfer(uploader net.Conn, downloader net.Conn) {
	SendString(uploader, "start", COMMANDSIZE)
	fileName, fileSizeData := getFileNameAndSize(uploader)
	sendFileNameAndSize(downloader, fileName, fileSizeData)
	fileSize, _ := strconv.ParseInt(fileSizeData, 10, 64)

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(downloader, uploader, (fileSize - receivedBytes))
			//Get the filler bytes
			io.CopyN(downloader, uploader, (receivedBytes+BUFFERSIZE)-fileSize)
			break
		}
		io.CopyN(downloader, uploader, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
}

//RelayHTTPTransfer - TODO
func RelayHTTPTransfer(uploader net.Conn, w http.ResponseWriter) {
	SendString(uploader, "start", COMMANDSIZE)

	fileName, fileSizeData := getFileNameAndSize(uploader)
	fileSize, _ := strconv.ParseInt(fileSizeData, 10, 64)
	w.Header().Add("Content-Disposition", "inline; filename=\""+fileName+"\"")

	//TODO - move this to a function
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(w, uploader, (fileSize - receivedBytes))
			//Get the filler bytes
			io.CopyN(w, uploader, (receivedBytes+BUFFERSIZE)-fileSize)
			break
		}
		io.CopyN(w, uploader, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
}

func getFileNameAndSize(connection net.Conn) (fileName string, fileSize string) {
	var err error
	fileName, err = ReadString(connection, COMMANDSIZE)
	CheckError(err)
	fileSizeData, err := ReadString(connection, COMMANDSIZE)
	CheckError(err)
	return fileName, fileSizeData
}

func sendFileNameAndSize(connection net.Conn, fileName string, fileSize string) {
	SendString(connection, fileName, COMMANDSIZE)
	SendString(connection, fileSize, COMMANDSIZE)
}

//DownloadFile - TODO
func DownloadFile(connection net.Conn) string {
	fileName, fileSizeData := getFileNameAndSize(connection)
	fileSize, _ := strconv.ParseInt(fileSizeData, 10, 64)

	bar := pb.StartNew(int(fileSize))

	newFile, err := os.Create(fileName)
	CheckError(err)

	defer newFile.Close()

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			bar.Set(int(fileSize))
			bar.FinishPrint("Download completed")
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
		bar.Set(int(receivedBytes))
	}
	return fileName
}

//SendFile - TODO
func SendFile(connection net.Conn, filePath string) {
	file, err := os.Open(filePath)
	CheckFatalError(err)
	fileInfo, err := file.Stat()
	CheckFatalError(err)

	fileSizeData := fileInfo.Size()
	fileSize := strconv.FormatInt(fileSizeData, 10)
	fileName := fileInfo.Name()
	bar := pb.StartNew(int(fileSizeData))

	sendFileNameAndSize(connection, fileName, fileSize)

	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			bar.Set(int(fileSizeData))
			bar.FinishPrint("Transfer complete")
			break
		}
		bar.Add(BUFFERSIZE)
		connection.Write(sendBuffer)
	}
	return
}

//Debug - TODO
func Debug(msg string) {
	if DebugEnabled {
		log.Println(msg)
	}
}
