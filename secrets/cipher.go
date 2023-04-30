package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
)

func main() {
	key := flag.String("key", "", "Key to encrypt with")
	value := flag.String("value", "", "Value to encrypt")

	flag.Parse()

	encryptedString, err := Encrypt(*key, *value)
	if err != nil {
		panic(err)
	}
	fmt.Println(encryptedString)
	decryptedString, err := Decrypt(*key, encryptedString)
	if err != nil {
		panic(err)
	}
	fmt.Println(decryptedString)
}

func HashKey(key string) []byte {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	return hasher.Sum(nil)
}

func Encrypt(key, data string) (string, error) {

	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	plaintext := []byte(data)
	block, err := aes.NewCipher(HashKey(key))
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	fmt.Printf("%x\n", ciphertext)
	encryptedHex := hex.EncodeToString(ciphertext)

	return encryptedHex, nil
}

func Decrypt(key, encryptedString string) (string, error) {

	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	ciphertext, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(HashKey(key))
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	fmt.Printf("%s", ciphertext)
	return string(ciphertext), nil

}
