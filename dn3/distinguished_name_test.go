package dn

import (
	"testing"
)

func TestQuoteChar(t *testing.T) {
	for _, v := range []string{`"`, `\`} {
		if quotechar([]rune(v)) != nil {
			t.Errorf("value found for %s", v)
		}
	}
}

func TestStringChar(t *testing.T) {
	for _, v := range []string{`"`, `\`, `,`, `=`, `+`, `<`, `>`, `#`, `;`} {
		if stringchar([]rune(v)) != nil {
			t.Errorf("value found for %s", v)
		}
	}
}
