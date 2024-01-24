package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

func GenerateKeyPairBase64(bits int) (pubBase64 string, privBase64 string, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	pubBytes := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
	if err != nil {
		return "", "", err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	pubBase64 = base64.StdEncoding.EncodeToString(pubBytes)
	privBase64 = base64.StdEncoding.EncodeToString(privBytes)
	return pubBase64, privBase64, nil
}

func EncryptAesJson(data interface{}, priv string) (string, error) {
	privateKeyBytes, _ := base64.StdEncoding.DecodeString(priv)
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	cipherBytes, err := EncryptAes(jsonMsg, privateKeyBytes)
	if err != nil {
		return "", err
	}
	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)
	return cipherText, err
}

func DecodeAesJson(cipherText string, priv string, dest interface{}) error {
	priv64, _ := base64.StdEncoding.DecodeString(priv)
	cipher64, _ := base64.StdEncoding.DecodeString(cipherText)
	text, err := DecryptAes(cipher64, priv64)
	if err != nil {
		return err
	}
	err = json.Unmarshal(text, dest)
	if err != nil {
		return err
	}
	return nil
}

func EncryptAes(plainText []byte, key []byte) ([]byte, error) {
	hash := sha256.Sum256(key)
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, err
}

func DecryptAes(cipherText []byte, key []byte) ([]byte, error) {
	hash := sha256.Sum256(key)
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}
