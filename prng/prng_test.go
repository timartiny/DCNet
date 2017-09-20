package prng

import (
    "testing"
    "fmt"
)

func TestSetup(t *testing.T){

    fmt.Println("** [Test] [Setup]\n")
    prng := Prng{}
    var err error
    k32 := make([]byte, 32)

    err = prng.Setup(k32)
    t.Log("32 byte key err =", err)

    if len(prng.pt) == 0 {
        fmt.Println("Setup is not working correctly")
    }

    fmt.Println("\n** ----------------------------------\n\n")
}

func TestIncrement(t *testing.T){
    fmt.Println("** [Test] [Increment]\n")

    prng := Prng{}
    k32 := make([]byte, 32)
    prng.Setup(k32)

    for i := 0; i < 10; i++ {
        fmt.Printf("%d - [% x]\n", i,  prng.pt)
        prng.increment()
    }

    for j := 0; j < 250; j++ {
        prng.increment()
    }
    fmt.Printf("260 - [% x]\n", prng.pt)

    for k := 0; k < 65279; k++ {
        prng.increment()
    }
    fmt.Printf("65539 - [% x]\n", prng.pt)

    fmt.Println("\n** ----------------------------------\n\n")
}

func TestGetBytes_short(t *testing.T){

    fmt.Println("** [Test] [GetBytes] - Short\n")
    prng := Prng{}
    k32 := make([]byte, 32)
    prng.Setup(k32)


    bytes := prng.GetBytes(10)
    fmt.Printf("[% x]\n", bytes)
    fmt.Printf("Bytes used - %d\n", prng.used)

    fmt.Println("\n** ----------------------------------\n\n")
}


func TestGetBytes_Long(t *testing.T){

    fmt.Println("** [Test] [GetBytes] - Long\n")
    prng := Prng{}
    k32 := make([]byte, 32)
    prng.Setup(k32)


    bytes := prng.GetBytes(100)
    fmt.Printf("[% x]\n", bytes)
    fmt.Printf("Bytes used - %d\n", prng.used)

    fmt.Println("\n** ----------------------------------\n\n")
}
