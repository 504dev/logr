package log

import (
	"encoding/base64"
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/models/dashboard"
	"testing"
)

func TestFinder(t *testing.T) {
	dash := dashboard.GetById(1)
	fmt.Println("Dashboard:", dash)
	publicKeyBytes, _ := base64.StdEncoding.DecodeString(dash.PublicKey)
	privateKeyBytes, _ := base64.StdEncoding.DecodeString(dash.PrivateKey)
	publicKey, _ := cipher.BytesToPublicKey(publicKeyBytes)
	privateKey, _ := cipher.BytesToPrivateKey(privateKeyBytes)
	publicKeyBytes, _ = cipher.PublicKeyToBytes(publicKey)
	privateKeyBytes = cipher.PrivateKeyToBytes(privateKey)
	fmt.Println("PublicKey:", base64.StdEncoding.EncodeToString(publicKeyBytes))
	fmt.Println("PrivateKey:", base64.StdEncoding.EncodeToString(privateKeyBytes))
}
