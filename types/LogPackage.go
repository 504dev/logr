package types

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	*Log       `json:"log"`
	*Metr      `json:"metr"`
}

func (lp *LogPackage) DecryptLog(priv string) error {
	log := Log{}
	err := log.Decrypt(lp.CipherText, priv)
	if err != nil {
		return err
	}
	lp.Log = &log
	return nil
}
