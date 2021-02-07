package mail

import (
	"bytes"
	"log"
)

/*
=============
= TODO/BUGS =
=============
- Split headers into lines at most 75 characters long
- More strict Q_Encode (recursive calls?)
*/

var crlf = []byte{13, 10}

const hexCharacters = "0123456789ABCDEF"

func isPrintableASCIIByte(b byte) bool {
	return b > 31 && b < 127
}

func isPrintableASCIIString(s string) bool {
	ascii := true
	for _, b := range []byte(s) {
		ascii = ascii && isPrintableASCIIByte(b)
	}
	return ascii
}

func isPrintableASCIICRLFByte(b byte) bool {
	return (b > 31 && b < 127) || b == 13 || b == 10
}

func isPrintableASCIICRLFString(s string) bool {
	ascii := true
	for _, b := range []byte(s) {
		ascii = ascii && isPrintableASCIICRLFByte(b)
	}
	return ascii
}

// The "Q" endcoding for message headers, defined in RFC2047 http://tools.ietf.org/html/rfc2047#section-4.2
func Q_Encode(s string) string {
	if isPrintableASCIIString(s) {
		return s
	}
	buf := bytes.NewBufferString("=?UTF-8?Q?")
	for _, b := range []byte(s) {
		if b == 32 { // Space is encoded with _
			buf.WriteByte(95)
		} else if isPrintableASCIIByte(b) && b != 95 && b != 61 && b != 63 {
			buf.WriteByte(b)
		} else {
			buf.WriteByte(61)
			buf.WriteByte(hexCharacters[b>>4])
			buf.WriteByte(hexCharacters[b&15])
		}
	}
	buf.WriteString("?=")
	o := buf.String()
	if len(o) > 76 {
		log.Println("To long Q-encoding")
	}
	return o
}

// Quoted-Printable encoding, http://en.wikipedia.org/wiki/Quoted-printable
func QP_Encode(s string) string {
	var buf bytes.Buffer
	lastNewline := 0
	for _, b := range []byte(s) {
		if lastNewline > 70 { // Break long lines
			lastNewline = 0
			buf.WriteByte(61)
			buf.Write(crlf)
		}
		if b == 10 { // New lines are terminated with CRLF
			lastNewline = 0
			buf.Write(crlf)
		} else if isPrintableASCIIByte(b) && b != 61 { // Normal character, except =
			lastNewline++
			buf.WriteByte(b)
		} else { // Escapable character
			lastNewline += 3
			buf.WriteByte(61)
			buf.WriteByte(hexCharacters[b>>4])
			buf.WriteByte(hexCharacters[b&15])
		}
	}
	return buf.String()
}
