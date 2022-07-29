package gstr

import (
	"github.com/Ravior/gserver/internal/utils"
)

// IsNumeric tests whether the given string s is numeric.
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}
