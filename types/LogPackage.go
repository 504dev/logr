package types

type LogPackage struct {
	DashId        int    `json:"dash_id"`
	PublicKey     string `json:"public_key"`
	CipherLog     string `json:"cipher_log"`
	CipherCounter string `json:"cipher_counter"`
	*Log
	*Counter
}

func (lp *LogPackage) DecryptLog(priv string) error {
	log := Log{}
	err := log.Decrypt(lp.CipherLog, priv)
	if err != nil {
		return err
	}
	lp.Log = &log
	return nil
}
