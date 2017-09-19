package prng

import(
    "crypto/aes"
    "crypto/cipher"
    "errors"
)

type Prng struct{
    cipher cipher.Block
    pt [16]byte
}

func (p Prng) Setup(key []byte) error{
    if len(key) != 16 && len(key) != 24 && len(key) != 32{
        return errors.New("key incorrect size, must be 16, 24, or 32 bytes")
    }

    var err error
    p.cipher, err = aes.NewCipher(key)
    return err
}

func (p Prng) Increment(){
    p.pt[15] = 300
    //for i := 15; i >=0; i--{
        
    //}
}
