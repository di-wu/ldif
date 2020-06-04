package dn

import (
	"testing"
)

func TestStringChar(t *testing.T) {
	for _, v := range []string{`"`, `\`, `,`, `=`, `+`, `<`, `>`, `#`, `;`} {
		if stringchar([]rune(v)) != nil {
			t.Errorf("value found for %s", string(v))
		}
	}
}

func TestExample(t *testing.T) {
	s := []rune(`CN=Marshall T. Rose, O=Dover Beach Consulting, L=Santa Clara, ST=California, C=US`)
	dn := name(s).Best()

	list := dn.GetSubNodes(`name-component`)
	if l := len(list); l != 5 {
		t.Errorf("should have 5 components, got %d", l)
	}

	if v := list[0].GetSubNode(`key`).Value; string(v) != "CN" {
		t.Errorf("expected \"CN\", got %s", string(v))
	}
	if v := list[0].GetSubNode(`string`).Value; string(v) != "Marshall T. Rose" {
		t.Errorf("expected \"Marshall T. Rose\", got %s", string(v))
	}

	if v := list[1].GetSubNode(`key`).Value; string(v) != "O" {
		t.Errorf("expected \"O\", got %s", string(v))
	}
	if v := list[1].GetSubNode(`string`).Value; string(v) != "Dover Beach Consulting" {
		t.Errorf("expected \"Dover Beach Consulting\", got %s", string(v))
	}

	if v := list[2].GetSubNode(`key`).Value; string(v) != "L" {
		t.Errorf("expected \"L\", got %s", string(v))
	}
	if v := list[2].GetSubNode(`string`).Value; string(v) != "Santa Clara" {
		t.Errorf("expected \"Santa Clara\", got %s", string(v))
	}

	if v := list[3].GetSubNode(`key`).Value; string(v) != "ST" {
		t.Errorf("expected \"ST\", got %s", string(v))
	}
	if v := list[3].GetSubNode(`string`).Value; string(v) != "California" {
		t.Errorf("expected \"California\", got %s", string(v))
	}

	if v := list[4].GetSubNode(`key`).Value; string(v) != "C" {
		t.Errorf("expected \"C\", got %s", string(v))
	}
	if v := list[4].GetSubNode(`string`).Value; string(v) != "US" {
		t.Errorf("expected \"US\", got %s", string(v))
	}
}

func TestInformalDefinition(t *testing.T) {
	for _, str := range []string{
		// folded
		`CN=Steve Kille,
O=ISODE Consortium, C=GB`,
		// multi-column format
		`CN=Steve Kille,
O=ISODE Consortium,
C=GB`,
		`CN=Christian Huitema, O=INRIA, C=FR`,
		// semicolon (";")
		`CN=Christian Huitema; O=INRIA; C=FR`,
		// different attribute types
		`CN=James Hacker,
L=Basingstoke,
O=Widget Inc,
C=GB`,
		// multi-valued relative distinguished name
		`OU=Sales + CN=J. Smith, O=Widget Inc., C=US`,
		//  quoting of a comma
		`CN=L. Eagle, O="Sue, Grabbit and Runn", C=GB`,
		`CN=L. Eagle, O=Sue\, Grabbit and Runn, C=GB`,
	} {
		if name([]rune(str)) == nil {
			t.Errorf("could not parse string: %s", str)
		}
	}
}

func TestExamples(t *testing.T) {
	for _, str := range []string{
		`CN=Marshall T. Rose, O=Dover Beach Consulting, L=Santa Clara, ST=California, C=US`,
		`CN=FTAM Service, CN=Bells, OU=Computer Science, O=University College London, C=GB`,
		`CN=Markus Kuhn, O=University of Erlangen, C=DE`,
		`CN=Steve Kille,
O=ISODE Consortium,
C=GB`,
		`CN=Steve Kille ,

O =   ISODE Consortium,
C=GB`,
		`CN=Steve Kille, O=ISODE Consortium, C=GB`,
	} {
		if name([]rune(str)) == nil {
			t.Errorf("could not parse string: %s", str)
		}
	}
}
