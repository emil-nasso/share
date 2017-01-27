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
	connection net.Conn
}

//Read - Implementation of the io.Reader interface
func (s *ShareConnection) Read(p []byte) (n int, err error) {
	return s.connection.Read(p)
}

//Write - Implementation of the io.Writer interface
func (s *ShareConnection) Write(p []byte) (n int, err error) {
	return s.connection.Write(p)
}

//Connect - TODO
func Connect(hostname, port string) *ShareConnection {
	connection, err := net.Dial("tcp", hostname+":"+port)
	lib.CheckFatalError(err)

	s := &ShareConnection{connection: connection}
	s.NegotiateVersion()
	return s
}

//New - Create a ShareConnection based on an existing net.Conn
func New(connection net.Conn) *ShareConnection {
	return &ShareConnection{
		connection: connection,
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
	if s.connection == nil {
		panic("Not connected")
	}
}

//Disconnect - TODO
func (s *ShareConnection) Disconnect() {
	s.connection.Close()
	s.connection = nil
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
	_, err := s.connection.Read(buffer)
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
	paddedString := fillString(str, length)
	lib.Debug("Sent string: " + str)
	s.connection.Write([]byte(paddedString))
}

//FillString - TODO There has to be a better way
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
