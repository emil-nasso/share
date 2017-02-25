package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/emil-nasso/share/conn"
	"github.com/emil-nasso/share/lib"
)

func relayFileTransfer(uploader *conn.ShareConnection, downloader *conn.ShareConnection) {
	fileName, fileSizeStr, fileSize := initializeFileTransfer(uploader)
	downloader.SendFileNameAndSize(fileName, fileSizeStr)
	copy(uploader, downloader, fileSize)
}

//RelayHTTPTransfer - TODO
func relayHTTPTransfer(uploader *conn.ShareConnection, w http.ResponseWriter) {
	fileName, _, fileSize := initializeFileTransfer(uploader)
	w.Header().Add("Content-Disposition", "inline; filename=\""+fileName+"\"")
	copy(uploader, w, fileSize)
}

func initializeFileTransfer(uploader *conn.ShareConnection) (fileName, sizeStr string, sizeInt int64) {
	uploader.SendString("start", lib.COMMANDSIZE)
	fileName, fileSizeString := uploader.GetFileNameAndSize()
	fileSize, _ := strconv.ParseInt(fileSizeString, 10, 64)
	return fileName, fileSizeString, fileSize
}

func copy(src io.Reader, dst io.Writer, fileSize int64) {
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) <= lib.BUFFERSIZE {
			io.CopyN(dst, src, (fileSize - receivedBytes))
			//Get the filler bytes
			io.CopyN(dst, src, (receivedBytes+lib.BUFFERSIZE)-fileSize)
			break
		}
		io.CopyN(dst, src, lib.BUFFERSIZE)
		receivedBytes += lib.BUFFERSIZE
	}
}
