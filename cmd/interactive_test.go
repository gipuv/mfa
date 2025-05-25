package cmd

import (
	"fmt"
	"testing"
)

func TestInteractive(t *testing.T) {
	tests := []string{"abc", "JBSWY3DPEHPK3PXP", "123", "A1B2C3", "MZXW6YTBOI======"}

	for _, s := range tests {
		fmt.Printf("%q valid? %v\n", s, isValidSecret(s))
	}
}
