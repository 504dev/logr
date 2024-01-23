package types

import (
	"bytes"
	_types "github.com/504dev/logr-go-client/types"
)

type LogPackageRow []*_types.LogPackage

func (chunks LogPackageRow) Joined() (complete bool, joined *_types.LogPackage) {
	complete = true
	ciphered := chunks[0].CipherLog != ""
	var buffer bytes.Buffer
	for _, lp := range chunks {
		if lp == nil {
			complete = false
			break
		}
		if ciphered {
			buffer.WriteString(lp.CipherLog)
		} else {
			buffer.WriteString(lp.PlainLog)
		}
	}

	if !complete {
		return false, nil
	}

	clone := *chunks[0]
	if ciphered {
		clone.CipherLog = buffer.String()
	} else {
		clone.PlainLog = buffer.String()
	}

	return true, &clone
}
