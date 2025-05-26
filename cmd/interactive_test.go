package cmd

import (
	"fmt"
	"testing"

	"github.com/gipuv/mfa/util"
)

func TestInteractive(t *testing.T) {
	tests := []string{"abc", "JBSWY3DPEHPK3PXP", "123", "A1B2C3", "MZXW6YTBOI======"}

	for _, s := range tests {
		// 去除空格，转大写
		s = util.PadBase32Secret(s)
		fmt.Printf("%q valid? %v\n", s, isValidSecret(s))
	}
}
