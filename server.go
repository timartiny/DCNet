package main

import (
	// "bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"./message"
)

var port = "0.0.0.0:9001"

var connections []net.Conn

func exit(conn net.Conn){
	for i, v := range connections{
		if v == conn{
			connections[i] = connections[len(connections)-1]
			connections = connections[:len(connections)-1]
			fmt.Println("Connection destroyed")
			return
		}
	}
}

// echo receives messages, and sends them to all connected clients.
// currently does no clean up on connections.
func echo(conn net.Conn) {
	// r := bufio.NewReader(conn)
	loop := true
	for loop {
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		switch err {
		case nil:
			break
		case io.EOF:
			return
			break
		default:
			fmt.Println("ERROR", err)
			return
			break
		}
		protoData := new(message.Message)
		err = proto.Unmarshal(data[0:n], protoData)
		if *protoData.Type == 3{
			exit(conn)
			loop = false
		}
		// fmt.Printf("Received Message %s, sending back\n", protoData.String())
		for _,v := range connections{
			v.Write(data[0:n])
		}
	}
}

// starts server.
func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("ERROR", err)
			continue
		}
		fmt.Println("Connection Received")
		connections = append(connections, conn)
		go echo(conn)
	}
}
