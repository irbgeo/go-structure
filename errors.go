package structure

import (
	"errors"
)

var (
	NeedPrtTypeErr = errors.New("needs '*struct' type for sample")
	immutableErr   = errors.New("src is immutable")
)
