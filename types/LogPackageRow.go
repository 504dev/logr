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

	ciphered := row[0].CipherLog != ""
	var buffer bytes.Buffer

	for _, lp := range row {
		if ciphered {
			buffer.WriteString(lp.CipherLog)
		} else {
			buffer.WriteString(lp.PlainLog)
		}
	}

	clone := *row[0]
	if ciphered {
		clone.CipherLog = buffer.String()
	} else {
		clone.PlainLog = buffer.String()
	}

	return true, &clone
}
