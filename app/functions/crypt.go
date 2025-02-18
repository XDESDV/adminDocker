package functions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

var passphrase string

// GetPassphrase Récupération de la passphrase
func GetPassphrase() string {
	return passphrase
}

// SetPassphrase alimentation de la passphrase
func SetPassphrase(newpassphrase string) {
	passphrase = newpassphrase
}

// Called by Encrypt to encrypt datas.
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt Crypte une chaine
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	var (
		err        error
		block      cipher.Block
		gcm        cipher.AEAD
		nonce      []byte
		ciphertext []byte
	)
	if block, err = aes.NewCipher([]byte(createHash(passphrase))); err == nil {
		if gcm, err = cipher.NewGCM(block); err == nil {
			nonce = make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
				ciphertext = gcm.Seal(nonce, nonce, data, nil)
			}
		}
	}
	return ciphertext, err
}

// Decrypt Decrypte une chaine
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	var (
		plaintext []byte
		err       error
	)
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return plaintext, err
	}
	return plaintext, err
}
