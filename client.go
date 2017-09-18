package main

import (
	"bufio"
	"bytes"
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
var curve ecdh.ECDH
var sharedKeys map[string][]byte

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
		protoData := new(message.Message)
		err = proto.Unmarshal(data[0:n], protoData)
		if err != nil{
			log.Fatal("error unmarshalling: ", err)
		}
		// fmt.Println("Received:", protoData.String())

		// now we need to parse received message
		go parseMessage(protoData, conn)
	}
}

func parseMessage(msg *message.Message, conn net.Conn){
	if *msg.Type == 0{
		//received a key, check if it is ours, if it is not, check if we already
		// have the shared key.
		if bytes.Equal(msg.Data, curve.Marshal(pubKey)){
			//this is our key
			fmt.Println("this is my Key")
			return
		}

		// not our key, generate a secret, see if we already have it.
		otherKey, fail := curve.Unmarshal(msg.Data)
		if fail != true {
			fmt.Println("Sent key couldn't unmarshal")
			return
		}

		_, ok := sharedKeys[string(msg.Data)]
		if ok{ 
			// already have this key
			fmt.Println("i have this key")
			return
		}
		fmt.Println("don't have key")
		
		// don't have this key
		if sharedKeys == nil{
			sharedKeys = make(map[string][]byte)
		}
		newSharedKey, err := curve.GenerateSharedSecret(privKey, otherKey)
		if err != nil {
			fmt.Println("couldn't generate shared secret, err=",err)
			return
		}

		sharedKeys[string(msg.Data)] = newSharedKey

		fmt.Println("added new key!")
		// now your key has to be sent in response
		conn.SetReadDeadline(time.Now())
		sendKey(conn)
		conn.SetReadDeadline(time.Time{})
	}
}

func genKeys() {
	privKey, pubKey, _  = curve.GenerateKey(rand.Reader)
}

func sendKey(conn net.Conn){
	m := &message.Message{
		Data: curve.Marshal(pubKey),
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

	curve = ecdh.NewCurve25519ECDH() 
	genKeys()
	sendKey(conn)

	go receive(conn)
	userInput := bufio.NewReader(os.Stdin)
	for {
		// fmt.Print("Message: ")
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
