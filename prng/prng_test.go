package prng

import (
    "testing"
    "fmt"
)



func TestGetBytes_short(t *testing.T){

    fmt.Println("** [Test] [GetBytes] - Short\n")

    k32 := make([]byte, 32)
    nonce := make([]byte, 16)


    bytes, err := GetBytes(k32, 10, nonce)

    if err != nil {
        t.Log(err)
    }

    fmt.Printf("[% x]\n", bytes)

    fmt.Println("\n** ----------------------------------\n\n")
}


func TestGetBytes_Long(t *testing.T){

    fmt.Println("** [Test] [GetBytes] - Long\n")

    k32 := make([]byte, 32)
    nonce := make([]byte, 16)


    bytes, err := GetBytes(k32, 100, nonce)

    if err != nil {
        t.Log(err)
    }

    fmt.Printf("[% x]\n", bytes)

    fmt.Println("\n** ----------------------------------\n\n")
}
