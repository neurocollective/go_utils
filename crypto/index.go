package crypto

import (
  "crypto/aes"
  "crypto/rand"
  "crypto/cipher"
  "encoding/base64"
  "crypto/rsa"
  "crypto/sha512"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "os"
)

func GenerateRandomBytes(n int) ([]byte, error) {
  b := make([]byte, n)
  _, err := rand.Read(b)
  if err != nil {
    return nil, err
  }

  return b, nil
}

// This should be in an env file in production

func Base64Encode(b []byte) string {
  return base64.StdEncoding.EncodeToString(b)
}

func Base64Decode(s string) []byte {
  data, err := base64.StdEncoding.DecodeString(s)
  if err != nil {
    panic(err)
  }
  return data
}

// Encrypt method is to encrypt or hide any classified text
func EncryptString(text, key string, iv []byte) (string, error) {
  block, err := aes.NewCipher([]byte(key))
  if err != nil {
    return "", err
  }
  plainText := []byte(text)
  cfb := cipher.NewCFBEncrypter(block, iv)
  cipherText := make([]byte, len(plainText))
  cfb.XORKeyStream(cipherText, plainText)
  return Base64Encode(cipherText), nil
}

func Encrypt(text, key, iv []byte) (string, error) {
  block, err := aes.NewCipher(key)
  if err != nil {
    return "", err
  }
  cfb := cipher.NewCFBEncrypter(block, iv)
  cipherText := make([]byte, len(text))
  cfb.XORKeyStream(cipherText, text)
  return Base64Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func DecryptString(text, key string, iv []byte) (string, error) {
  block, err := aes.NewCipher([]byte(key))
  if err != nil {
    return "", err
  }
  cipherText := Base64Decode(text)
  cfb := cipher.NewCFBDecrypter(block, iv)
  plainText := make([]byte, len(cipherText))
  cfb.XORKeyStream(plainText, cipherText)
  return string(plainText), nil
}

func Decrypt(text, key, iv []byte) ([]byte, error) {
  block, err := aes.NewCipher([]byte(key))
  if err != nil {
    return []byte{}, err
  }
  cipherText := text
  cfb := cipher.NewCFBDecrypter(block, iv)
  plainText := make([]byte, len(cipherText))
  cfb.XORKeyStream(plainText, cipherText)
  return plainText, nil
}

/*
true?

comment here: 
https://stackoverflow.com/questions/5583379/what-is-the-limit-to-the-amount-of-data-that-can-be-encrypted-with-rsa

"Technically an RSA public key is nothing more than 2 integers 
(a really big modulus and an exponent). You can leave out everything else 
associated with it and still encrypt something the private key holder can decrypt.""

`rsa.PrivateKey` struct in `crypto/rsa` seems to confirm this.
*/

/* 
  how big can a message be?
  `const getMessageSize = (keyBitSize, hashSize) => {
    return (Math.floor(keyBitSize/8) - (2 * Math.ceil(hashSize/8)) - 2);
  };`
*/

/*

https://stackoverflow.com/questions/13378815/base64-length-calculation

plaintext to base64. If max base64 encryped transmission is 

So you need `4*(n/3)` chars to represent n bytes, and this needs to be 
rounded up to a multiple of 4.

*/

func RSATest() error {
  // begin generate private key
  fmt.Println("generating key...")


  // 20480 & sha512 allows message size of 2430 bytes 
  privateKey, err := rsa.GenerateKey(rand.Reader, 10240)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error generating RSA key: %s", err)
    return err
  }

  fmt.Println("done generating key.")

  thePKCS8PrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error marshalling RSA private key: %s", err)
    return err
  }

  pemBytes := pem.EncodeToMemory(&pem.Block{
    Type:  "PRIVATE KEY",
    Bytes: thePKCS8PrivateKey,
  })
  fmt.Println("pem:", string(pemBytes))

  // done generating private key

  raw := "send reinforcements, we're going to advance asdfasdfasdfasdfadsfadsfassdfaaaaaaaaaaaaaaaaaadddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
  secretMessage := []byte(raw)
  label := []byte{}

  // crypto/rand.Reader is a good source of entropy for randomizing the
  // encryption function.
  rng := rand.Reader

  fmt.Println("Encrypting...")
  ciphertext, err := rsa.EncryptOAEP(sha512.New(), rng, &privateKey.PublicKey, secretMessage, label)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
    return err
  }

  // Since encryption is a randomized function, ciphertext will be
  // different each time.
  fmt.Printf("done encrypting. Ciphertext: %x\n", ciphertext)

  // type OAEPOptions struct {
  //   // Hash is the hash function that will be used when generating the mask.
  //   Hash crypto.Hash

  //   // MGFHash is the hash function used for MGF1.
  //   // If zero, Hash is used instead.
  //   MGFHash crypto.Hash

  //   // Label is an arbitrary byte string that must be equal to the value
  //   // used when encrypting.
  //   Label []byte
  // }

  // options := rsa.OAEPOptions{
  //   Hash: sha256.New(),
  //   Label: []byte{},
  // }

  // plainText, err := privateKey.Decrypt(rng, ciphertext, options)

  fmt.Println("Decrypting...")
  plainText, err := rsa.DecryptOAEP(sha512.New(), rng, privateKey, ciphertext, []byte{})

  if err != nil {
    return err
  }

  fmt.Println("done decrypting. plaintext:", string(plainText))

  return nil
}
