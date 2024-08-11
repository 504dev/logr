package types

import _types "github.com/504dev/logr-go-client/types"

type LogPackageMeta struct {
	*_types.LogPackage
	Protocol string
	Size     int
}
