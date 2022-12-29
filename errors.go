package structure

import (
	"errors"
)

var (
	NeedPrtTypeErr = errors.New("needs '*struct' type")
	immutableErr   = errors.New("src is immutable")
)
