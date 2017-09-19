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

// receive handles the receiving of messages.
// it puts the received data into a protobuf and passes it to a parsing function.
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

// parseMessage will take a message and connection and will determine how to
// handle a message based on the type of the message. Some messages require 
// immediate responses, and those are given.
func parseMessage(msg *message.Message, conn net.Conn){
	if *msg.Type == 0{
		newKey := parseKey(msg.Data)

		if newKey{
			// now your key has to be sent in response
			conn.SetReadDeadline(time.Now())
			sendKey(conn)
			conn.SetReadDeadline(time.Time{})
		}
	}
}

// parseKey will parse given data to see if it is a key we already have, if not
// it will add it to our global map, and generate a shared secret. 
// Returns whether the key is new.
func parseKey(data []byte) bool{
	//received a key, check if it is ours, if it is not, check if we already
	// have the shared key.
	if bytes.Equal(data, curve.Marshal(pubKey)){
		//this is our key
		fmt.Println("this is my Key")
		return false
	}

	// not our key, generate a secret, see if we already have it.
	otherKey, fail := curve.Unmarshal(data)
	if fail != true {
		fmt.Println("Sent key couldn't unmarshal")
		return false
	}

	_, ok := sharedKeys[string(data)]
	if ok{ 
		// already have this key
		fmt.Println("i have this key")
		return false
	}
	fmt.Println("don't have key")
	
	// don't have this key
	if sharedKeys == nil{
		sharedKeys = make(map[string][]byte)
	}

	newSharedKey, err := curve.GenerateSharedSecret(privKey, otherKey)
	if err != nil {
		fmt.Println("couldn't generate shared secret, err=",err)
		return false
	}

	sharedKeys[string(data)] = newSharedKey

	fmt.Println("added new key!")
	return true;
}

// genKeys generates keys based on curve ed25519, fills in global values.
func genKeys() {
	privKey, pubKey, _  = curve.GenerateKey(rand.Reader)
}

// sendKey will handle sending out the public Keys, should be generalized for
// sending a message
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

// main client method, establishes connection to server, allows user to type
// messages to be sent.
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
