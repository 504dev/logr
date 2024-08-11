package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

func GenerateKeyPairBase64(bits int) (pubBase64 string, privBase64 string, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	pubBytes := sha256.Sum256(x509.MarshalPKCS1PublicKey(&priv.PublicKey))
	privBytes := sha256.Sum256(x509.MarshalPKCS1PrivateKey(priv))
	pubBase64 = base64.StdEncoding.EncodeToString(pubBytes[:])
	privBase64 = base64.StdEncoding.EncodeToString(privBytes[:])
	return pubBase64, privBase64, nil
}
