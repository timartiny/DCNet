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
		if err != nil {
			fmt.Println("err unmarshalling proto, err=",err)
		}
		if *protoData.Type == 4{
			exit(conn)
			loop = false
		}
		fmt.Printf("Received Message [% x], sending back\n", protoData.Data[:3])
		for _,v := range connections{
			n, err := v.Write(data[0:n])
            fmt.Printf(" [% x] bytes: [%s]\n", data[0:n], v.RemoteAddr())
			if err != nil{
				fmt.Println("err sending data, err=",err)
			}
		}
        fmt.Print("\n")
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
		fmt.Println("Connectfion Received")
		connections = append(connections, conn)
		go echo(conn)
	}
}
