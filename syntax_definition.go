package ldif

import . "github.com/di-wu/abnf"

func File(s []rune) *AST {
	return Alts(
		`ldif-file`,
		ldifContent,
		ldifChanges,
	)(s)
}

func ldifContent(s []rune) *AST {
	return Concat(
		`ldif-content`,
		versionSpec,
		Repeat1Inf(`1*(1*SEP ldif-attrval-record)`, Concat(
			`1*SEP ldif-attrval-record`,
			Repeat1Inf(`1*SEP`, sep),
			ldifAttrvalRecord,
		)),
	)(s)
}

func ldifChanges(s []rune) *AST {
	return Concat(
		`ldif-changes`,
		versionSpec,
		Repeat1Inf(`1*(1*SEP ldif-change-record)`, Concat(
			`1*SEP ldif-change-record`,
			Repeat1Inf(`1*SEP`, sep),
			ldifChangeRecord,
		)),
	)(s)
}

func ldifAttrvalRecord(s []rune) *AST {
	return Concat(
		`ldif-attrval-record`,
		dnSpec,
		sep,
		Repeat1Inf(`1*attrval-spec`, attrvalSpec),
	)(s)
}

func ldifChangeRecord(s []rune) *AST {
	return Concat(
		`ldif-change-record`,
		dnSpec,
		sep,
		Repeat0Inf(`*control`, control),
		changerecord,
	)(s)
}

func versionSpec(s []rune) *AST {
	return Concat(
		`version-spec`,
		String(`version:`, "version:", true),
		fill,
		versionNumber,
	)(s)
}

// MUST be "1" for the LDIF format described in this document.
func versionNumber(s []rune) *AST {
	return Repeat1Inf(`version-number`, digit)(s)
}

func dnSpec(s []rune) *AST {
	return Concat(
		`dn-spec`,
		String(`dn:`, "dn:", true),
		Alts(
			`FILL distinguishedName / ":" FILL base64-distinguishedName`,
			Concat(
				`FILL distinguishedName`,
				fill,
				distinguishedName,
			),
			Concat(
				`":" FILL base64-distinguishedName`,
				Rune(`:`, ':'),
				fill,
				base64DistinguishedName,
			),
		),
	)(s)
}

// a distinguished name, as defined in RFC2253
func distinguishedName(s []rune) *AST {
	return safeString(s)
}

// a distinguishedName which has been base64 encoded
func base64DistinguishedName(s []rune) *AST {
	return base64utf8String(s)
}

// a relative distinguished name, defined as <name-component> in RFC2253
func rdn(s []rune) *AST {
	return safeString(s)
}

// an rdn which has been base64 encoded
func base64Rdn(s []rune) *AST {
	return base64utf8String(s)
}

func control(s []rune) *AST {
	return Concat(
		`control`,
		String(`control:`, "control:", true),
		fill,
		ldapOid, // control type
		Repeat(`0*1(1*SPACE ("true" / "false"))`, 0, 1, Concat(
			`1*SPACE ("true" / "false")`,
			Repeat1Inf(`1*SPACE`, space),
			Alts(
				`"true" / "false"`,
				String(`true`, "true", true),
				String(`false`, "false", true),
			),
		)),                                         // criticality
		Repeat(`0*1(value-spec)`, 0, 1, valueSpec), // control value
		sep,
	)(s)
}

func ldapOid(s []rune) *AST {
	return Concat(
		`ldap-oid`,
		Repeat1Inf(`1*DIGIT`, digit),
		Repeat(`0*1("." 1*DIGIT)`, 0, 1, Concat(
			`0*1("." 1*DIGIT)`,
			Rune(`.`, '.'),
			Repeat1Inf(`1*DIGIT`, digit),
		)),
	)(s)
}

func attrvalSpec(s []rune) *AST {
	return Concat(
		`attrval-spec`,
		attributeDescription,
		valueSpec,
		sep,
	)(s)
}

func valueSpec(s []rune) *AST {
	return Concat(
		`value-spec`,
		Rune(`:`, ':'),
		Alts(
			`FILL 0*1(SAFE-STRING) / ":" FILL (BASE64-STRING) / "<" FILL url`,
			Concat(
				`FILL 0*1(SAFE-STRING)`,
				fill,
				Repeat(`0*1(SAFE-STRING)`, 0, 1, safeString),
			),
			Concat(
				`":" FILL (BASE64-STRING)`,
				Rune(`:`, ':'),
				fill,
				base64String,
			),
			Concat(
				`"<" FILL url)`,
				Rune(`<`, '<'),
				fill,
				url,
			),
		),
	)(s)
}

func url(s []rune) *AST {
	return Concat(`url`)(s) // TODO
}

func attributeDescription(s []rune) *AST {
	return Concat(
		`AttributeDescription`,
		attributeType,
		Optional(`[";" options]`, Concat(
			`";" options`,
			Rune(`;`, ';'),
			options,
		)),
	)(s)
}

func attributeType(s []rune) *AST {
	return Alts(
		`AttributeType`,
		ldapOid,
		Concat(
			`ALPHA *(attr-type-chars)`,
			alpha,
			Repeat0Inf(`*(attr-type-chars)`, attrTypeChars),
		),
	)(s)
}

func options(s []rune) *AST {
	return Alts(
		`options`,
		option,
		Concat(
			`option ";" options`,
			option,
			Rune(`;`, ';'),
			options,
		),
	)(s)
}

func option(s []rune) *AST {
	return Repeat1Inf(`option`, optChar)(s)
}

func attrTypeChars(s []rune) *AST {
	return Alts(
		`attr-type-chars`,
		alpha,
		digit,
		Rune(`-`, '-'),
	)(s)
}

func optChar(s []rune) *AST {
	return attrTypeChars(s)
}

func changerecord(s []rune) *AST {
	return Concat(
		`changerecord`,
		String(`changetype:`, "changetype:", true),
		fill,
		Alts(
			`change-add / change-delete / change-modify / change-moddn`,
			changeAdd,
			changeDelete,
			changeModify,
			changeModdn,
		),
	)(s)
}

func changeAdd(s []rune) *AST {
	return Concat(
		`change-add`,
		String(`add`, "add", true),
		sep,
		Repeat1Inf(`1*attrval-spec`, attrvalSpec),
	)(s)
}

func changeDelete(s []rune) *AST {
	return Concat(
		`change-delete`,
		String(`delete`, "delete", true),
		sep,
	)(s)
}

func changeModdn(s []rune) *AST {
	return Concat(
		`change-moddn`,
		Alts(
			`"modrdn" / "moddn"`,
			String(`modrdn`, "modrdn", true),
			String(`moddn`, "moddn", true),
		),
		sep,
		String(`newrdn:`, "newrdn:", true),
		Alts(
			`FILL rdn / ":" FILL base64-rdn`,
			Concat(
				`FILL rdn`,
				fill,
				rdn,
			),
			Concat(
				`":" FILL base64-rdn`,
				Rune(`:`, ':'),
				fill,
				base64Rdn,
			),
		),
		sep,
		String(`deleteoldrdn:`, "deleteoldrdn:", true),
		fill,
		Alts(
			`"0" / "1"`,
			Rune(`0`, '0'),
			Rune(`1`, '1'),
		),
		sep,
		Repeat(`0*1("newsuperior:" (FILL distinguishedName / ":" FILL base64-distinguishedName) SEP)`, 0, 1, Concat(
			`"newsuperior:" (FILL distinguishedName / ":" FILL base64-distinguishedName) SEP`,
			String(`newsuperior:`, "newsuperior:", true),
			Alts(
				`FILL distinguishedName / ":" FILL base64-distinguishedName`,
				Concat(
					`FILL distinguishedName`,
					fill,
					distinguishedName,
				),
				Concat(
					`":" FILL base64-distinguishedName`,
					Rune(`:`, ':'),
					fill,
					base64DistinguishedName,
				),
			),
			sep,
		)),
	)(s)
}

func changeModify(s []rune) *AST {
	return Concat(
		`change-modify`,
		String(`modify`, "modify", true),
		sep,
		Repeat0Inf(`*mod-spec`, modSpec),
	)(s)
}

func modSpec(s []rune) *AST {
	return Concat(
		`mod-spec`,
		Alts(
			`"add:" / "delete:" / "replace:"`,
			String(`add:`, "add:", true),
			String(`delete:`, "delete:", true),
			String(`replace:`, "replace:", true),
		),
		fill,
		attributeDescription,
		sep,
		Repeat0Inf(`*attrval-spec`, attrvalSpec),
		Rune(`-`, '-'),
		sep,
	)(s)
}

var (
	space = Rune(`SPACE`, '\x20') // ASCII SP, space
	fill  = Repeat0Inf(`FILL`, space)
	sep   = Alts(`SEP`, Concat(`CR LF`, cr, lf), lf)
	cr    = Rune(`CR`, '\x0D') // ASCII CR, carriage return
	lf    = Rune(`LF`, '\x0A') // ASCII LF, line feed
	alpha = Alts(
		`ALPHA`,
		Range(`%x41-5A`, '\x41', '\x5A'), // A-Z
		Range(`%x61-7A`, '\x61', '\x7A'), // a-z
	)
	digit = Range(`DIGIT`, '\x30', '\x39') // 0-9
	utf81 = Range(`UTF8-1`, '\x80', '\xBF')
	utf82 = Concat(`UTF8-2`, Range(`%xC0-DF`, '\xC0', '\xDF'), utf81)
	utf83 = Concat(`UTF8-3`, Range(`%xE0-EF`, '\xE0', '\xEF'), RepeatN(`2UTF8-1`, 2, utf81))
	utf84 = Concat(`UTF8-4`, Range(`%xF0-F7`, '\xF0', '\xF7'), RepeatN(`3UTF8-1`, 3, utf81))
	utf85 = Concat(`UTF8-5`, Range(`%xF8-FB`, '\xF8', '\xFB'), RepeatN(`4UTF8-1`, 4, utf81))
	utf86 = Concat(`UTF8-6`, Range(`%xFC-FD`, '\xFC', '\xFD'), RepeatN(`5UTF8-1`, 5, utf81))
	// any value <= 127 decimal except NUL, LF and CR
	safeChar = Alts(
		`SAFE-CHAR`,
		Range(`%x01-09`, '\x01', '\x09'),
		Range(`%x0B-0C`, '\x0B', '\x0C'),
		Range(`%x0E-7F`, '\x0E', '\x7F'),
	)
	// any value <= 127 except NUL, LF, CR, SPACE, colon (":", ASCII 58 decimal) and less-than ("<" , ASCII 60 decimal)
	safeInitChar = Alts(
		`SAFE-INIT-CHAR`,
		Range(`%x01-09`, '\x01', '\x09'),
		Range(`%x0B-0C`, '\x0B', '\x0C'),
		Range(`%x0E-1F`, '\x0E', '\x1F'),
		Range(`%x21-39`, '\x21', '\x39'),
		Rune(`%x3B`, '\x3B'),
		Range(`%x3D-7F`, '\x3D', '\x7F'),
	)
	safeString = Optional(`SAFE-STRING`, Concat(
		`SAFE-INIT-CHAR *SAFE-CHAR`,
		safeInitChar,
		Repeat0Inf(`*SAFE-CHAR`, safeChar),
	))
	utf8Char = Alts(
		`UTF8-CHAR`,
		safeChar,
		utf82,
		utf83,
		utf84,
		utf85,
		utf86,
	)
	utf8String = Repeat0Inf(`*UTF8-CHAR`, utf8Char)
	// MUST be the base64 encoding of a UTF8-STRING
	base64utf8String = base64String
	base64Char       = Alts(
		`BASE64-CHAR`,
		Rune(`%x2B`, '\x2B'),             // +
		Rune(`%x2F`, '\x2F'),             // /
		Range(`%x30-39`, '\x30', '\x39'), // 0-9
		Rune(`%x3D`, '\x3D'),             // =
		Range(`%x41-5A`, '\x41', '\x5A'), // A-Z
		Range(`%x61-7A`, '\x61', '\x7A'), // a-z
	)
	base64String = Optional(`BASE64-STRING`, Repeat0Inf(`*(BASE64-CHAR)`, base64Char))
)
