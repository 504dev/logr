package types

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	*Log       `json:"log"`
	*Metr      `json:"metr"`
}
