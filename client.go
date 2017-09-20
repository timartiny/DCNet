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
	"os/signal"
	// "time"

	"github.com/golang/protobuf/proto"
	"./go-ecdh"
	"./message"
)

var port = "0.0.0.0:9001"
var pubKey crypto.PublicKey
var privKey crypto.PrivateKey
var curve ecdh.ECDH
var sharedKeys map[string][]byte
var conn net.Conn

// receive handles the receiving of messages.
// it puts the received data into a protobuf and passes it to a parsing function.
func receive(){
	// response := bufio.NewReader(conn)
	for {
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		switch err {
		case nil:
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
		go parseMessage(protoData)
	}
}

// parseMessage will take a message and connection and will determine how to
// handle a message based on the type of the message. Some messages require 
// immediate responses, and those are given.
func parseMessage(msg *message.Message){
	if *msg.Type == 0{
		//received a key
		newKey := parseKey(msg.Data)

		if newKey{
			// now your key has to be sent in response
			// conn.SetReadDeadline(time.Now())
			m := &message.Message{
				Data: curve.Marshal(pubKey),
				Type: proto.Int32(0),
			}
			sendMessage(m)
			// conn.SetReadDeadline(time.Time{})
		}
	}else if *msg.Type == 1{
		//received a new message
		fmt.Printf("Received: %s\n", msg.Data)
	}else if *msg.Type == 3{
		//received a disconnect message
		delete(sharedKeys, string(msg.Data))
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
		// fmt.Println("this is my Key")
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
		// fmt.Println("i have this key")
		return false
	}
	// fmt.Println("don't have key")
	
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

	// fmt.Println("added new key!")
	return true;
}

// genKeys generates keys based on curve ed25519, fills in global values.
func genKeys() {
	curve = ecdh.NewCurve25519ECDH() 
	privKey, pubKey, _  = curve.GenerateKey(rand.Reader)
}

// sendMessage handles the sending of messages to the server, it will Marshal
// the protobuf and then write it to the server.
func sendMessage(m *message.Message){
	data, err := proto.Marshal(m)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	_, err = conn.Write(data)
	if err != nil{
		log.Fatal("sending error: ", err)
	}
	// fmt.Printf("Sent %d bytes\n", n)
}

func cleanup(){
	m:= &message.Message{
		Data: curve.Marshal(pubKey),
		Type: proto.Int32(3),
	}
	sendMessage(m)
}

// main client method, establishes connection to server, allows user to type
// messages to be sent.
func main() {
	// set up the connection
	var err error
	conn, err = net.Dial("tcp", port)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	// start receiving messages
	go receive()

	// we have two types of input, a typed message and also a ctrl+c signal
	// we create a channel for the signal first, if a user types ctrl+c it will
	// cause the anonymous function to call cleanup() and inform the main loop 
	// that it should exit
	sigChan := make(chan os.Signal, 1)
	sigHandled := make(chan bool)
	signal.Notify(sigChan, os.Interrupt)
	go func(){
		<- sigChan
		cleanup()
		sigHandled <- true
	}()

	// set up our keys
	genKeys()

	// send our keys to the server to broadcast to all users
	m := &message.Message{
		Data: curve.Marshal(pubKey),
		Type: proto.Int32(0),
	}
	sendMessage(m)

	// we now make a channel for our user input, since we have to listen for
	// two types of input we need two channels each running concurrently.
	// once the input has been gathered it will fill channel with input or 
	// handle error.
	inpChan := make(chan []byte)
	go func(inpChan chan []byte){
		for {
			fmt.Print("Message: ")
			userInput := bufio.NewReader(os.Stdin)
			userLine, err := userInput.ReadBytes(byte('\n'))
				switch err {
				case nil:
					inpChan <- userLine
				case io.EOF:
					os.Exit(0)
				default:
					fmt.Println("ERROR", err)
					os.Exit(1)
				}
		}
	}(inpChan)

	// main loop, waits for channels to be filled, if input is filled it will 
	// send the message, if the signal channel is filled it will handle cleanup
	// and then exit
	for {
		select{
		case inp :=<- inpChan:
			m = &message.Message{
				Data: inp,
				Type: proto.Int32(1),
			}
			sendMessage(m)
		case <- sigHandled:
			cleanup()
			os.Exit(0)
		}
	}
}
