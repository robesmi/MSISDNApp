package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/robesmi/MSISDNApp/model/errs"
)


// Encrypts a email using AE256 method and returns an encrypted string encoded in base 64...or an error
func EncryptEmailAes256(key []byte, message string) (string, error) {

	c, cErr := aes.NewCipher(key)
	if cErr != nil{
		return "", errs.NewEncryptionError(cErr.Error())
	}

	gcm, gcmErr := cipher.NewGCM(c)
	if gcmErr != nil{
		return "", errs.NewEncryptionError(cErr.Error())
	}

	// Uses the first 12 bytes of the message(email) as its nonce
	// So we can reliably get the same output needed for database lookup
	nonce := make([]byte, gcm.NonceSize())
	copy(nonce[:], []byte(message))


	result := gcm.Seal(nonce, nonce, []byte(message), nil)
	
	encodedResult := base64.RawStdEncoding.EncodeToString(result)

	return encodedResult, nil
}

// Decrypts a Aes256 email and returns the original string.... or an error?
func DecryptEmailAes256(key []byte, ciphertext string) (string, error){

	decodedCiphertext, decErr := base64.RawStdEncoding.DecodeString(ciphertext)
	if decErr != nil{
		return "", errs.NewEncryptionError(decErr.Error())
	}

	c, cErr := aes.NewCipher(key)
	if cErr != nil {
		return "", errs.NewEncryptionError(cErr.Error())
	}

	gcm, gcmErr := cipher.NewGCM(c)
	if gcmErr != nil{
		return "", errs.NewEncryptionError(cErr.Error())
	}

	nonceSize := gcm.NonceSize()
	if len(decodedCiphertext) < nonceSize{
		return "", errs.NewEncryptionError("Nonce size is less than the gcm nonce size")
	}

	nonce, cipher := decodedCiphertext[:nonceSize], decodedCiphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(cipher), nil)
	if err != nil{
		return "", errs.NewEncryptionError(err.Error())
	}

	return string(plaintext), nil

}