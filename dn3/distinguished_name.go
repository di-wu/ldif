package dn

// RFC 2253: 3. Parsing a String back to a Distinguished Name

import (
	"unicode/utf8"

	. "github.com/elimity-com/abnf/operators"
)

func distinguishedName(s []rune) Alternatives {
	return Optional(
		`distinguishedName`,
		name,
	)(s)
}

func name(s []rune) Alternatives {
	return Concat(
		`name`,
		nameComponent,
		Repeat0Inf(
			`*("," name-component)`,
			Concat(
				`"," name-component`,
				Rune(`,`, ','),
				nameComponent,
			),
		),
	)(s)
}

func nameComponent(s []rune) Alternatives {
	return Concat(
		`name-component`,
		attributeTypeAndValue,
		Repeat0Inf(
			`*("+" attributeTypeAndValue)`,
			Concat(
				`"+" attributeTypeAndValue`,
				Rune(`+`, '+'),
				attributeTypeAndValue,
			),
		),
	)(s)
}

func attributeTypeAndValue(s []rune) Alternatives {
	return Concat(
		`attributeTypeAndValue`,
		attributeType,
		Rune(`=`, '='),
		attributeValue,
	)(s)
}

func attributeType(s []rune) Alternatives {
	return Alts(
		`attributeType`,
		Concat(
			`ALPHA 1*keychar`,
			alpha,
			Repeat1Inf(
				`1*keychar`,
				keychar,
			),
		),
		oid,
	)(s)
}

func keychar(s []rune) Alternatives {
	return Alts(
		`keyChar`,
		alpha,
		digit,
		Rune(`-`, '-'),
	)(s)
}

func oid(s []rune) Alternatives {
	return Concat(
		`oid`,
		Repeat1Inf(`1*DIGIT`, digit),
		Repeat0Inf(`*("." 1*DIGIT)`, Concat(
			`"." 1*DIGIT`,
			Rune(`.`, '.'),
			Repeat1Inf(`1*DIGIT`, digit),
		)),
	)(s)
}

func attributeValue(s []rune) Alternatives {
	return str(s)
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
			`"#" hexstring`,
			Rune(`#`, '#'),
			hexstring,
		),
		Concat(
			`QUOTATION *(quotechar / pair) QUOTATION`,
			quotation,
			Repeat0Inf(`*(quotechar / pair)`, Alts(
				`quotechar / pair`,
				quotechar,
				pair,
			)),
			quotation,
		),
	)(s)
}

func quotechar(s []rune) Alternatives {
	return Alts(
		`quotechar`,
		Range(`0-33`, 0, 33),
		// except quotation
		Range(`35-91`, 35, 91),
		// except "\"
		Range(`93-`, 93, utf8.MaxRune),
	)(s)
}

func special(s []rune) Alternatives {
	return Alts(
		`special`,
		Rune(`,`, ','),
		Rune(`=`, '='),
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
			`special / "\" / QUOTATION / hexpair`,
			special,
			Rune(`\`, 92),
			quotation,
			hexpair,
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

func hexstring(s []rune) Alternatives {
	return Repeat1Inf(`hexstring`, hexpair)(s)
}

func hexpair(s []rune) Alternatives {
	return Concat(`hexpair`, hexchar, hexchar)(s)
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
	alpha = Alts(
		`ALPHA`,
		Range(`%x41-5A`, 65, 90),  // 65-90
		Range(`%x61-7A`, 97, 122), // 97-122
	)
	digit     = Range(`DIGIT`, 48, 57) // 45-57
	quotation = Rune(`QUOTATION`, 34)  // 34
)
