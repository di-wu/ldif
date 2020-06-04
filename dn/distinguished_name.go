package dn

import (
	"unicode/utf8"

	. "github.com/elimity-com/abnf/operators"
)

// RFC 1779: 2.3 Formal definition

func name(s []rune) Alternatives {
	return Alts(
		`name`,
		Concat(
			`name-component [spaced-separator]`,
			nameComponent,
			Optional(`[spaced-separator]`, spacedSeparator),
		),
		Concat(
			`name-component spaced-separator name`,
			nameComponent,
			spacedSeparator,
			name,
		),
	)(s)
}

func spacedSeparator(s []rune) Alternatives {
	return Concat(
		`spaced-separator`,
		optionalSpace,
		separator,
		optionalSpace,
	)(s)
}

func separator(s []rune) Alternatives {
	return Alts(
		`separator`,
		Rune(`,`, ','),
		Rune(`;`, ';'),
	)(s)
}

func optionalSpace(s []rune) Alternatives {
	return Concat(
		`optional-space`,
		Optional(`[cr]`, cr),
		Repeat0Inf(`*(" ")`, Rune(` `, ' ')),
	)(s)
}

func nameComponent(s []rune) Alternatives {
	return Alts(
		`name-component`,
		attribute,
		Concat(
			`attribute optional-space "+" optional-space name-component`,
			attribute,
			optionalSpace,
			Rune(`+`, '+'),
			optionalSpace,
			nameComponent,
		),
	)(s)
}

func attribute(s []rune) Alternatives {
	return Alts(
		`attribute`,
		str,
		Concat(
			`key optional-space "=" optional-space string`,
			key,
			optionalSpace,
			Rune(`=`, '='),
			optionalSpace,
			str,
		),
	)(s)
}

func key(s []rune) Alternatives {
	return Alts(
		`key`,
		Repeat1Inf(`1*(keychar)`, keychar),
		Concat(
			`"OID." oid`,
			StringCS(`OID.`, "OID."),
			oid,
		),
		Concat(
			`"oid." oid`,
			StringCS(`oid.`, "oid."),
			oid,
		),
	)(s)
}

func keychar(s []rune) Alternatives {
	return Alts(
		`keyChar`,
		alpha,
		digit,
		Rune(` `, ' '),
	)(s)
}

func oid(s []rune) Alternatives {
	return Alts(
		`oid`,
		digitstring,
		Concat(
			`digitstring "." oid`,
			digitstring,
			Rune(`.`, '.'),
			oid,
		),
	)(s)
}

func digitstring(s []rune) Alternatives {
	return Repeat1Inf(`digitstring`, digit)(s)
}

// string
func str(s []rune) Alternatives {
	return Alts(
		`string`,
		Repeat0Inf(`*(stringchar / pair)`, Alts(
			`stringchar / pair`,
			stringchar,
			pair,
		)),
		Concat(
			`""" *(stringchar / special / pair) """`,
			Rune(`"`, '"'),
			Repeat0Inf(`*(stringchar / special / pair)`, Alts(
				`stringchar / special / pair`,
				stringchar,
				special,
				pair,
			)),
			Rune(`"`, '"'),
		),
		Concat(
			`"#" <hex>`,
			Rune(`#`, '#'),
			hex,
		),
	)(s)
}

func special(s []rune) Alternatives {
	return Alts(
		`special`,
		Rune(`,`, ','),
		Rune(`=`, '='),
		cr,
		Rune(`+`, '+'),
		Rune(`<`, '<'),
		Rune(`>`, '>'),
		Rune(`#`, '#'),
		Rune(`;`, ';'),
	)(s)
}

func pair(s []rune) Alternatives {
	return Concat(
		`pair`,
		Rune(`\`, 92), // "\"
		Alts(
			`special / "\" / """`,
			special,
			Rune(`\`, 92),
			Rune(`"`, '"'),
		),
	)(s)
}

func stringchar(s []rune) Alternatives {
	return Alts(
		`stringchar`,
		Range(`0-33`, 0, 33),
		// except quotation, and "#"
		Range(`36-42`, 36, 42),
		// except "," and "+"
		Range(`45-58`, 45, 58),
		// except ";", "<", "=" and ">"
		Range(`63-91`, 63, 91),
		// except "\"
		Range(`93-`, 93, utf8.MaxRune),
	)(s)
}

func hex(s []rune) Alternatives {
	return Repeat(`hex`, 2, -1, hexchar)(s)
}

func hexchar(s []rune) Alternatives {
	return Alts(`hexchar`, digit,
		Rune(`A`, 'A'), Rune(`B`, 'B'), Rune(`C`, 'C'),
		Rune(`D`, 'D'), Rune(`E`, 'E'), Rune(`F`, 'F'),
		Rune(`a`, 'a'), Rune(`b`, 'b'), Rune(`c`, 'c'),
		Rune(`d`, 'd'), Rune(`e`, 'e'), Rune(`f`, 'f'),
	)(s)
}

var (
	cr    = Rune(`CR`, 13)
	alpha = Alts(
		`ALPHA`,
		Range(`%x41-5A`, 65, 90),  // 65-90
		Range(`%x61-7A`, 97, 122), // 97-122
	)
	digit = Range(`DIGIT`, 48, 57) // 45-57
)
