package log

import (
	"fmt"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/mysql"
	"testing"
)

func TestCrypt(t *testing.T) {
	var err error
	config.Init()
	mysql.Init()
	//priv, pub, _ := cipher.GenerateKeyPair(256)
	//pubBytes, _ := cipher.PublicKeyToBytes(pub)
	//privBytes := cipher.PrivateKeyToBytes(priv)
	//fmt.Println(base64.StdEncoding.EncodeToString(pubBytes))
	//fmt.Println(base64.StdEncoding.EncodeToString(privBytes))
	dash, err := dashboard.GetById(1)
	fmt.Println(dash, err)
	if err != nil {
		panic(err)
	}
	logpack := LogPackage{
		PublicKey: dash.PublicKey,
		Log:       GetLast(),
	}
	fmt.Println("Create:", logpack.Log, logpack.CipherText)
	err = logpack.EncryptLog()
	logpack.Log = nil
	fmt.Println("Encrypt:", logpack.Log, logpack.CipherText, err)
	err = logpack.DecryptLog()
	fmt.Println("Decrypt:", logpack.Log, logpack.CipherText, err)
}
