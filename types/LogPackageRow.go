package types

import (
	"bytes"
	_types "github.com/504dev/logr-go-client/types"
)

type LogPackageRow []*_types.LogPackage

func (row LogPackageRow) Joined() (complete bool, joined *_types.LogPackage) {
	complete = true
	for _, lp := range row {
		if lp == nil {
			complete = false
			break
		}
	}

	if !complete {
		return false, nil
	}

	ciphered := row[0].CipherLog != nil
	var buffer bytes.Buffer

	for _, lp := range row {
		if ciphered {
			buffer.Write(lp.CipherLog)
		} else {
			buffer.Write(lp.PlainLog)
		}
	}

	clone := *row[0]
	if ciphered {
		clone.CipherLog = buffer.Bytes()
	} else {
		clone.PlainLog = buffer.Bytes()
	}

	return true, &clone
}
