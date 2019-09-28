package log

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/models/dashboard"
	"testing"
)

func TestFinder(t *testing.T) {
	//priv, pub, _ := cipher.GenerateKeyPair(4096)
	//pubBytes, _ := cipher.PublicKeyToBytes(pub)
	//privBytes := cipher.PrivateKeyToBytes(priv)
	//fmt.Println(base64.StdEncoding.EncodeToString(pubBytes))
	//fmt.Println(base64.StdEncoding.EncodeToString(privBytes))
	dash := dashboard.GetById(1)
	privateKeyBytes, _ := base64.StdEncoding.DecodeString(dash.PrivateKey)
	privateKey, _ := cipher.BytesToPrivateKey(privateKeyBytes)
	fmt.Println("Dashboard:", dash)
	logitem := GetLast()
	jsonMsg, _ := json.Marshal(logitem)
	fmt.Println("Json:", string(jsonMsg))
	cipherBytes, err := cipher.EncryptWithPublicKey(jsonMsg, &privateKey.PublicKey)
	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)
	fmt.Println("CipherText:", cipherText, err)
	sigBytes, err := cipher.SignMessage(jsonMsg, privateKey)
	sig := base64.StdEncoding.EncodeToString(sigBytes)
	fmt.Println("Sig:", sig, err)
	logpack := LogPackage{
		PublicKey:  dash.PublicKey,
		CipherText: cipherText,
		Signature:  sig,
	}
	fmt.Println("LogPackage:", logpack)
	err = logpack.Decrypt()
	fmt.Println("Log:", logpack.Log, err)
}
