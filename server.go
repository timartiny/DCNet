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

func echo(conn net.Conn) {
	// r := bufio.NewReader(conn)
	for {
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
		// fmt.Printf("Received Message %s, sending back\n", protoData.String())
		for _,v := range connections{
			v.Write(data[0:n])
		}
	}
}

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
