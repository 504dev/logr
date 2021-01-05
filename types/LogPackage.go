package types

type LogPackage struct {
	DashId      int    `json:"dash_id,omitempty"`
	PublicKey   string `json:"public_key"`
	CipherLog   string `json:"cipher_log,omitempty"`
	CipherCount string `json:"cipher_count,omitempty"`
	*Log        `json:"log,omitempty"`
	*Count      `json:"count,omitempty"`
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

func (lp *LogPackage) DecryptCount(priv string) error {
	count := Count{}
	err := count.Decrypt(lp.CipherCount, priv)
	if err != nil {
		return err
	}
	lp.Count = &count
	return nil
}
