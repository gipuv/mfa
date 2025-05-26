package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func TestInteractive(t *testing.T) {
	tests := []string{"abc", "JBSWY3DPEHPK3PXP", "123", "A1B2C3", "MZXW6YTBOI======"}

	for _, s := range tests {
		// 去除空格，转大写
		// totp.padBase32Secret
		s = strings.ToUpper(strings.ReplaceAll(s, " ", ""))
		missing := len(s) % 8
		if missing != 0 {
			s += strings.Repeat("=", 8-missing) // 补齐 '='
		}
		fmt.Printf("%q valid? %v\n", s, isValidSecret(s))
	}
}
