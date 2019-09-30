package log

import (
	"fmt"
	"github.com/504dev/kidlog/models/dashboard"
	"testing"
)

func TestCrypt(t *testing.T) {
	//priv, pub, _ := cipher.GenerateKeyPair(256)
	//pubBytes, _ := cipher.PublicKeyToBytes(pub)
	//privBytes := cipher.PrivateKeyToBytes(priv)
	//fmt.Println(base64.StdEncoding.EncodeToString(pubBytes))
	//fmt.Println(base64.StdEncoding.EncodeToString(privBytes))
	dash := dashboard.GetById(1)
	logpack := LogPackage{
		PublicKey: dash.PublicKey,
		Log:       GetLast(),
	}
	fmt.Println("Create:", logpack.Log, logpack.CipherText)
	var err error
	err = logpack.EncryptLog()
	logpack.Log = nil
	fmt.Println("Encrypt:", logpack.Log, logpack.CipherText, err)
	err = logpack.DecryptLog()
	fmt.Println("Decrypt:", logpack.Log, logpack.CipherText, err)
}
