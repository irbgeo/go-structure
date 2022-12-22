package structtagbuilder

import (
	"errors"
)

var (
	needPrtTypeErr    = errors.New("needs '*struct' type for sample")
	acceptableTypeErr = errors.New("the type of dst is not acceptable")
	immutableErr      = errors.New("src is immutable")
)
