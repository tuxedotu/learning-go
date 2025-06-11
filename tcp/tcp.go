package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/tuxedotu/learning-go/playground-db"
)

const SrvAddr = "localhost"
const SrvPort = "42069"

func client() {
	// set msg based on arguments
	msg := "Hello Server!" // -- default
	if len(os.Args) > 1 {
		msg = strings.Join(os.Args[1:len(os.Args)], " ")
	}

	// Dial up connetion to srv
	connection, err := net.Dial("tcp", SrvAddr+":"+SrvPort)
	if err != nil {
		fmt.Println(err)
		fmt.Println("net.Dial broke")
		return
	}

	fmt.Println("Sending: ", msg)
	err = gob.NewEncoder(connection).Encode(msg)
}

func server() {
	// tcp socket setup (listener)
	listener, err := net.Listen("tcp", SrvAddr+":"+SrvPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	// db setup
	loggingDb, err := playDB.OpenTcpLogsDB()
	if err != nil {
		fmt.Println(err)
	}
	defer loggingDb.Close()

	// client handling
	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnectionOnServer(connection, loggingDb)
		// clientMsg := make(chan string)
		// go func() { clientMsg <- handleConnectionOnServer(connection, loggingDb) }()
		// fmt.Println(<-clientMsg)

	}
}

func handleConnectionOnServer(currentConn net.Conn, db *sql.DB) string {
	var msg string
	err := gob.NewDecoder(currentConn).Decode(&msg)

	if err != nil {
		fmt.Println(err)
	} else {
		playDB.InsertTcpLog(db, currentConn.RemoteAddr().String(), msg)
		fmt.Println("Received (and stored): ", msg)
	}

	currentConn.Close()
	return msg
}

// MAIN
func main() {
	server()
	// client()
}
