package prng

import(
    "crypto/aes"
    "crypto/cipher"
    "errors"
    // "fmt"
)



func GetBytes( key []byte, n int, nonce []byte) ([]byte, error) {

    if len(key) != 16 && len(key) != 24 && len(key) != 32{
        return nil, errors.New("key incorrect size, must be 16, 24, or 32 bytes")
    }

    cipherText := make([]byte, aes.BlockSize+ n )
    plainText := make([]byte, n)

    for i:=0; i < n; i++{ plainText[i] = 0 }

    // initilialize the cipher
    cipherBlock, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    cfbStream := cipher.NewCFBEncrypter(cipherBlock, nonce)
    cfbStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
    return cipherText[aes.BlockSize:], nil

}





