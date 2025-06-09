package tcp

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

const SrvAddr = "localhost"
const SrvPort = "9999"

func Client() {
	// set msg based on arguments
	msg := "Hello Server!"

	if len(os.Args) > 1 {
		msg = strings.Join(os.Args[1:len(os.Args)], " ")
	}

	// Dial up connetion to srv-address
	connection, err := net.Dial("tcp", SrvAddr+":"+SrvPort)
	if err != nil {
		fmt.Println(err)
		fmt.Println("net.Dial broke")
		return
	}

	fmt.Println("Sending: ", msg)
	err = gob.NewEncoder(connection).Encode(msg)
}

func Server() {
	// listen to a port
	listener, err := net.Listen("tcp", SrvAddr+":"+SrvPort)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleServerConnection(connection)
	}
}

func handleServerConnection(currentConn net.Conn) {
	var msg string
	err := gob.NewDecoder(currentConn).Decode(&msg)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Received: ", msg)
	}

	currentConn.Close()
}
