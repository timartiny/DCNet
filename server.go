package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

var port = "0.0.0.0:9001"

var connections []net.Conn

func echo(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes(byte('\n'))
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
		fmt.Printf("Received Message %s, sending back\n", string(line))
		for _,v := range connections{
			v.Write(line)
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
