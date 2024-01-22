package types

import (
	"bytes"
	_types "github.com/504dev/logr-go-client/types"
)

type LogPackageRow []*_types.LogPackage

func (chunks LogPackageRow) Joined() (complete bool, joined *_types.LogPackage) {
	complete = true
	var buffer bytes.Buffer
	for _, lp := range chunks {
		if lp == nil {
			complete = false
			break
		}
		buffer.WriteString(lp.CipherLog)
	}

	if !complete {
		return false, nil
	}

	clone := *chunks[0]
	clone.CipherLog = buffer.String()

	return true, &clone
}
