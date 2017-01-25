package lib

import "log"

//BUFFERSIZE - The size of the buffer for sending files. This amount of bytes will be
//sent per chunk.
const BUFFERSIZE = 1024

//COMMANDSIZE - The default size of the chucks for sending command and for the responses
const COMMANDSIZE = 64

//PROTOCOLVERSION - the current client/server communication protocol version
const PROTOCOLVERSION = "v1"

//DebugEnabled - Is debug enabled?
var DebugEnabled bool

//TODO: http://golangtutorials.blogspot.se/2011/06/inheritance-and-subclassing-in-go-or.html
// lägg ut send och read i en struct, en communicator, använd "arv" för att både servern och klienten ska ha tillgång till den

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

//Debug - TODO
func Debug(msg string) {
	if DebugEnabled {
		log.Println(msg)
	}
}
