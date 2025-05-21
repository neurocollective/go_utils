package main

import (
  "github.com/neurocollective/go_utils/crypto"
  "fmt"
  // "os"
)

func main() {

  err := crypto.RSATest()

  if err != nil {
    fmt.Println(err)
  }

  // iv, err := crypto.GenerateRandomBytes(16)

  // keyBytes, err := crypto.GenerateRandomBytes(32)

  // if err != nil {
  //   fmt.Println("GenerateRandomBytes error: ", err)
  // }

  // key := string(keyBytes)

  // if err != nil {
  //   fmt.Println("GenerateRandomBytes error: ", err)
  // }

  // stringToEncrypt := "shhh"

  // encrypted, err := crypto.EncryptString(stringToEncrypt, key, iv)

  // if err != nil {
  //   fmt.Println("error encrypting your classified text: ", err)
  // }
  // fmt.Println(encrypted)

  // decText, err := crypto.DecryptString(encrypted, key, iv)
  // if err != nil {
  //   fmt.Println("error decrypting your encrypted text: ", err)
  // }
  // fmt.Println(decText)
}

// func main() {
//   iv, err := crypto.GenerateRandomBytes(16)

//   keyBytes, err := crypto.GenerateRandomBytes(16)

//   if err != nil {
//     fmt.Println("GenerateRandomBytes error: ", err)
//   }

//   key := keyBytes

//   if err != nil {
//     fmt.Println("GenerateRandomBytes error: ", err)
//   }

//   zipBytes, err := os.ReadFile("test.zip")

//   if err != nil {
//     fmt.Println("read zip error: ", err)
//   }

//   encrypted, err := crypto.Encrypt(zipBytes, key, iv)

//   if err != nil {
//     fmt.Println("error encrypting your classified text: ", err)
//   }

//   encryptedBytes := crypto.Base64Decode(encrypted)

//   err = os.WriteFile("test.enc.zip", encryptedBytes, 0666)

//   if err != nil {
//     fmt.Println("write file error: ", err)
//   }

//   encryptedZipBytes, err := os.ReadFile("test.enc.zip")

//   if err != nil {
//     fmt.Println("read file error: ", err)
//   }

//   decryptedZipBytes, err := crypto.Decrypt(encryptedZipBytes, key, iv)
//   if err != nil {
//     fmt.Println("error decrypting your encrypted text: ", err)
//   }
//   err = os.WriteFile("test.dec.zip", decryptedZipBytes, 0666)
// }

// func main() {

//  cwd, getCwdError := os.Getwd()

//  if getCwdError != nil {
//    os.Exit(1)
//  }

//  log.Println("cwd:", cwd)

//  fields := []map[string]string{
//    map[string]string{"fieldName": "id", "type": "int"},
//    map[string]string{"fieldName": "name", "type": "string"},
//  }

//  config := generator.GenerationConfig{
//    fields,
//    "TestStruct",
//    cwd + "/generated",
//    cwd + "/generator/index.templ",
//    "generatorTest",
//  }

//  generator.Generate(config)
// }
