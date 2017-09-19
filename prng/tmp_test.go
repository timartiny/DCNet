package prng

import "testing"

func TestSetup(t *testing.T){
    prng := Prng{}
    err := prng.Setup([]byte("hello"))
    t.Log("hello key err =", err)
    k16 := make([]byte, 16)
    k24 := make([]byte, 24)
    k32 := make([]byte, 32)

    err = prng.Setup(k16)
    t.Log("16 byte key err =", err)
    err = prng.Setup(k24)
    t.Log("24 byte key err =", err)
    err = prng.Setup(k32)
    t.Log("32 byte key err =", err)
}

func TestIncrement(t *testing.T){
    prng := Prng{}
    k16 := make([]byte, 16)
    prng.Setup(k16)
}
