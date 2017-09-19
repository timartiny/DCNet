package main

import (
    "fmt"
    "math/rand"
    "unicode/utf8"
    "encoding/binary"
    "errors"
    "os"
)


func NewRNG(secret []byte) *rand.Rand {
    s := rand.NewSource( secretToSeed(secret) )
    return rand.New(s)
}


func secretToSeed(secret []byte) int64{
    return int64( binary.BigEndian.Uint64(secret[0:7]) )
}


func safeXOR(dst, a, b []byte) error {
    n := len(a)
    if n > len(b) {
        n = len(b)
    }
    for i:=0; i < n; i++ {
        dst[i] =  a[i] ^ b[i]
    }
    return nil
}


func getRandomBytes(length int, buf []byte, r *rand.Rand) error {
    n := r.Read(buf)
    if n != length{
        return errors.New(" -- [CreateMessage] Random number generator error")
    }
    return nil
}

func createMessage(message string, RNGs map[string]*rand.Rand) ([]byte, error) {
    bytes := []byte(message)
    messageLen := utf8.RuneCountInString(message)
    hold := make( []byte, messageLen )

    for k := range RNGs {
        _ = RNGs[k].Read(hold)

        err := safeXOR(bytes, bytes, hold)
        if err !=nil {
            return nil, err
        }
    }

    return bytes, nil
}


func main(){
    secret := make([]byte, 32)
    var RNGs map[string]*rand.Rand

    RNGs[string(secret)] = NewRNG(secret)

    message := "Hello there friend"
    fmt.Println(message)

    bytes, err := createMessage(message, RNGs)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for i :=0; i < len(message); i++ {
        fmt.Printf("%X ", bytes[i])
    }
    fmt.Println("\n")
}


