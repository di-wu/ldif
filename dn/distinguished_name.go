package dn

import (
	"unicode/utf8"

	. "github.com/di-wu/abnf"
)

// RFC 1779: 2.3 Formal definition

func name(s []rune) *AST {
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

func spacedSeparator(s []rune) *AST {
	return Concat(
		`spaced-separator`,
		optionalSpace,
		separator,
		optionalSpace,
	)(s)
}

func separator(s []rune) *AST {
	return Alts(
		`separator`,
		Rune(`,`, ','),
		Rune(`;`, ';'),
	)(s)
}

func optionalSpace(s []rune) *AST {
	return Concat(
		`optional-space`,
		Optional(`[cr]`, cr),
		Repeat0Inf(`*(" ")`, Rune(` `, ' ')),
	)(s)
}

func nameComponent(s []rune) *AST {
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

func attribute(s []rune) *AST {
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

func key(s []rune) *AST {
	return Alts(
		`key`,
		Repeat1Inf(`1*(keychar)`, keychar),
		Concat(
			`"OID." oid`,
			String(`OID.`, "OID.", true),
			oid,
		),
		Concat(
			`"oid." oid`,
			String(`oid.`, "oid.", true),
			oid,
		),
	)(s)
}

func keychar(s []rune) *AST {
	return Alts(
		`keyChar`,
		alpha,
		digit,
		Rune(` `, ' '),
	)(s)
}

func oid(s []rune) *AST {
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

func digitstring(s []rune) *AST {
	return Repeat1Inf(`digitstring`, digit)(s)
}

// string
func str(s []rune) *AST {
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

func special(s []rune) *AST {
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

func pair(s []rune) *AST {
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

func stringchar(s []rune) *AST {
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

func hex(s []rune) *AST {
	return Repeat(`hex`, 2, -1, hexchar)(s)
}

func hexchar(s []rune) *AST {
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
