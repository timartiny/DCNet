package prng

import(
    "crypto/aes"
    "crypto/cipher"
    "errors"
    // "fmt"
)

type Prng struct{
    cipher  cipher.Block
    pt      []byte        // Plaintext number to Encrypt for next ct block
    ct      []byte        // Cipher Text Block
    used    int           // [0-15] Number of bytes used in Cipher text 
}

func (p *Prng) Setup(key []byte) error{
    if len(key) != 16 && len(key) != 24 && len(key) != 32{
        return errors.New("key incorrect size, must be 16, 24, or 32 bytes")
    }

    var err error
    // initilialize the cipher
    p.cipher, err = aes.NewCipher(key)
    p.pt = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

    // Generate the first Cipher text
    p.generateBlock()

    return err
}

func (p *Prng) increment(){
    carry := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    p.pt[0] += 1
    if p.pt[0] == 0 { carry[0] = 1 }

    for i := 1; i <16; i++{
        if (p.pt[i] + byte(carry[i-1]) == 0) && (carry[i-1] ==1) { carry[i] = 1 }
        p.pt[i] += byte(carry[i-1])
    }

    if carry[15] != 0 {
        // WE HAVE CARRIED OVER THROUGH THE ENTIRE PLAINTEXT
    }
}

func (p *Prng) getByte() byte{

    if p.used == 16 {
        p.generateBlock()
    }

    // fmt.Printf("%d\n", p.used)
    // fmt.Printf("ct - [% x]\n", p.ct)
    // fmt.Printf("pt - [% x]\n", p.pt)

    b := p.ct[p.used]
    p.used++

    return b
}

func (p *Prng) GetBytes(n int) []byte {
    hold := make([]byte, n)

    // Fill in the message one block of ciphertext at a time.
    for i := 0; i < n; i++ {
        hold[i] = p.getByte()
    }

    return hold
}


func (p *Prng) generateBlock() {

    var ciphertext = make([]byte, aes.BlockSize+len(p.pt))
    iv := ciphertext[:aes.BlockSize]

    cfb := cipher.NewCFBEncrypter(p.cipher, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], p.pt)
    p.ct = ciphertext[16:]

    // fmt.Printf("New Cipher Text [% x]\n", p.ct)

    p.increment()
    p.used = 0
}





