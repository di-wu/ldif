# LDAP Data Interchange Format: LDIF
[RFC2849](https://tools.ietf.org/html/rfc2849)

## Definition
The LDIF format is used to convey directory information, or a description of a set of changes made to directory entries.
An LDIF file consists of a series of records separated by line separators.
A record consists of a sequence of lines describing a directory entry, 
or a sequence of lines describing a set of changes to a directory entry.
An LDIF file specifies a set of directory entries, or a set of changes to be applied to directory entries, but not both.

## Formal Syntax Definition of LDIF
[ABNF Parser](https://github.com/di-wu/abnf/)
