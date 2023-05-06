package Clients

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func EncryteFile() string {

	rawKey, _ := KeyGenerateSH256()
	key := []byte(rawKey[0:16]) // 16 bytes
	plaintext := []byte("apidatabasekey")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	// Encrypt the file
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Write the encrypted file
	f, err := os.Create("encrypted")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(ciphertext)
	if err != nil {
		panic(err)
	}

	return string(key)

}

func DecryptFile() string {

	rawKey := EncryteFile()

	key := []byte(rawKey) // 16 bytes

	ciphertext, err := ioutil.ReadFile("encrypted")
	if err != nil {
		panic(err.Error())
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		panic("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
