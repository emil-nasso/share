package conn

import (
	"errors"
	"io"
	"log"
	"net"
	"strings"

	"github.com/emil-nasso/share/lib"
)

//ShareConnection - TODO
type ShareConnection struct {
	Connection net.Conn
}

//Connect - TODO
func Connect(hostname, port string) *ShareConnection {
	connection, err := net.Dial("tcp", hostname+":"+port)
	lib.CheckFatalError(err)

	s := &ShareConnection{Connection: connection}
	s.NegotiateVersion()
	return s
}

//New - Create a ShareConnection based on an existing net.Conn
func New(connection net.Conn) *ShareConnection {
	return &ShareConnection{
		Connection: connection,
	}
}

//NegotiateVersion - TODO
func (s *ShareConnection) NegotiateVersion() {
	s.CheckConnected()
	s.SendString(lib.PROTOCOLVERSION, lib.COMMANDSIZE)
	serverVersion, err := s.ReadString(lib.COMMANDSIZE)
	if err != nil {
		log.Fatalln(err)
	}
	if serverVersion != lib.PROTOCOLVERSION {
		log.Fatalln("Client/server version mismatch. Server:", serverVersion, "Client:", lib.PROTOCOLVERSION)
	}
}

//CheckConnected - TODO
func (s *ShareConnection) CheckConnected() {
	if s.Connection == nil {
		panic("Not connected")
	}
}

//Disconnect - TODO
func (s *ShareConnection) Disconnect() {
	s.Connection.Close()
	s.Connection = nil
}

//GetFileNameAndSize - TODO
func (s *ShareConnection) GetFileNameAndSize() (fileName string, fileSize string) {
	var err error
	fileName, err = s.ReadString(lib.COMMANDSIZE)
	lib.CheckError(err)
	fileSizeData, err := s.ReadString(lib.COMMANDSIZE)
	lib.CheckError(err)
	return fileName, fileSizeData
}

//SendFileNameAndSize - TODO
func (s *ShareConnection) SendFileNameAndSize(fileName string, fileSize string) {
	s.SendString(fileName, lib.COMMANDSIZE)
	s.SendString(fileSize, lib.COMMANDSIZE)
}

//ReadString - Reads a string from the Connection
func (s *ShareConnection) ReadString(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := s.Connection.Read(buffer)
	if err == io.EOF {
		return "", errors.New("Connection closed")
	} else if err != nil {
		return "", err
	}
	response := strings.Trim(string(buffer), ":")
	lib.Debug("Read string:" + response)
	return response, nil
}

//SendString - sends a string to the Connection
func (s *ShareConnection) SendString(str string, length int) {
	paddedString := lib.FillString(str, length)
	lib.Debug("Sent string: " + str)
	s.Connection.Write([]byte(paddedString))
}
