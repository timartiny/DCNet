package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var port = "0.0.0.0:9001"

func receive(conn net.Conn){
	response := bufio.NewReader(conn)
	for {
		serverLine, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			fmt.Printf("received: %s\n", string(serverLine))
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Println("ERROR", err)
			os.Exit(2)
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	go receive(conn)
	userInput := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Message: ")
		userLine, err := userInput.ReadBytes(byte('\n'))
		switch err {
		case nil:
			conn.SetReadDeadline(time.Now())
			conn.Write(userLine)
			conn.SetReadDeadline(time.Time{})
			go receive(conn)
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Println("ERROR", err)
			os.Exit(1)
		}

		// serverLine, err := response.ReadBytes(byte('\n'))
		// switch err {
		// case nil:
		// 	fmt.Printf("received: %s\n", string(serverLine))
		// case io.EOF:
		// 	os.Exit(0)
		// default:
		// 	fmt.Println("ERROR", err)
		// 	os.Exit(2)
		// }
	}
}
