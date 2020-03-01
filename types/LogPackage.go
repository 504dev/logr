package types

type LogPackage struct {
	DashId     int    `json:"dash_id"`
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
