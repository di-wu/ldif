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

func TestExamples(t *testing.T) {
	for _, str := range []string{
		// 3 relative distinguished names
		`CN=Steve Kille,O=Isode Limited,C=GB`,
		// multi-valued
		`OU=Sales+CN=J. Smith,O=Widget Inc.,C=US`,
		// comma
		`CN=L. Eagle,O=Sue\, Grabbit and Runn,C=GB`,
		// carriage return
		`CN=Before` + "\x0D" + `After,O=Test,C=GB`,
		// BER encoding of an OCTET STRING containing two bytes
		`1.3.6.1.4.1.1466.0=#04024869,O=Test,C=GB`,
		// 5 letters: L, U, C WITH CARON, I, C WITH ACUTE
		`SN=Lu` + string(0xC48D) + `i` + string(0xC487),
	} {
		if name([]rune(str)) == nil {
			t.Errorf("could not parse string: %s", str)
		}
	}
}
