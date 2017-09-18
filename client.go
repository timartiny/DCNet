package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"./go-ecdh"
	"./message"
)

var port = "0.0.0.0:9001"
var pubKey crypto.PublicKey
var privKey crypto.PrivateKey

func receive(conn net.Conn){
	// response := bufio.NewReader(conn)
	for {
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		switch err {
		case nil:
			fmt.Printf("Received %d bytes\n", n)
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Println("ERROR", err)
			os.Exit(2)
		}
		fmt.Println("Decoding data into protobuf")
		protoData := new(message.Message)
		err = proto.Unmarshal(data[0:n], protoData)
		if err != nil{
			log.Fatal("error unmarshalling: ", err)
		}
		fmt.Println("Received:", protoData.String())
	}
}

func genKeys(e ecdh.ECDH) {
	privKey, pubKey, _  = e.GenerateKey(rand.Reader)
}

func sendKey(e ecdh.ECDH, conn net.Conn){
	m := &message.Message{
		Data: e.Marshal(pubKey),
		Type: proto.Int32(0),
	}

	data, err := proto.Marshal(m)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	n, err := conn.Write(data)
	if err != nil{
		log.Fatal("sending error: ", err)
	}
	fmt.Printf("Sent %d bytes for key\n", n)
}

func main() {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	var curve = ecdh.NewCurve25519ECDH() 
	genKeys(curve)
	sendKey(curve, conn)

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
	}
}
