package mail

import (
	"bytes"
	"encoding/base64"
	"log"
	"math/rand"
	"time"
)

// Basic MIME part, implmeneted by all subtypes
type MIMEPart interface {
	ContentType() string
	ContentTransferEncoding() string
	EncodedBody() []byte
	Headers() map[string]string
}

func Write(p MIMEPart, buf *bytes.Buffer) {
	if len(p.ContentType()) != 0 {
		buf.WriteString("Content-Type: ")
		buf.WriteString(p.ContentType())
		buf.Write(crlf)
	}
	if len(p.ContentTransferEncoding()) != 0 {
		buf.WriteString("Content-Transfer-Encoding: ")
		buf.WriteString(p.ContentTransferEncoding())
		buf.Write(crlf)
	}
	for key, value := range p.Headers() {
		buf.WriteString(key)
		buf.WriteString(": ")
		buf.WriteString(value)
		buf.Write(crlf)
	}
	buf.Write(crlf)
	buf.Write(p.EncodedBody())
	buf.Write(crlf)
	buf.Write(crlf)
}

// Basic text/plain
type MIMEPlain string

func (s MIMEPlain) ContentType() string {
	if isPrintableASCIICRLFString(string(s)) {
		return "text/plain; charset=\"us-ascii\""
	} else {
		return "text/plain; charset=\"utf-8\""
	}
}

func (s MIMEPlain) ContentTransferEncoding() string {
	if isPrintableASCIICRLFString(string(s)) {
		return "7bit"
	} else {
		return "quoted-printable"
	}
}

func (s MIMEPlain) EncodedBody() []byte {
	if isPrintableASCIICRLFString(string(s)) {
		return []byte(s)
	} else {
		return []byte(QP_Encode(string(s)))
	}
}
func (s MIMEPlain) Headers() map[string]string {
	return map[string]string{}
}

// Basic text/html
type MIMEHTML string

func (m MIMEHTML) ContentType() string {
	if isPrintableASCIICRLFString(string(m)) {
		return "text/html; charset=\"us-ascii\""
	} else {
		return "text/html; charset=\"utf-8\""
	}
}

func (m MIMEHTML) ContentTransferEncoding() string {
	if isPrintableASCIICRLFString(string(m)) {
		return "7bit"
	} else {
		return "quoted-printable"
	}
}

func (m MIMEHTML) EncodedBody() []byte {
	if isPrintableASCIICRLFString(string(m)) {
		return []byte(m)
	} else {
		return []byte(QP_Encode(string(m)))
	}
}

func (m MIMEHTML) Headers() map[string]string {
	return map[string]string{}
}

// Basic binary data
type MIMEBinary struct {
	Type              string
	Body              []byte
	AdditionalHeaders map[string]string
}

func (b MIMEBinary) ContentType() string {
	return b.Type
}

func (b MIMEBinary) ContentTransferEncoding() string {
	return "base64"
}

func (b MIMEBinary) EncodedBody() []byte {
	encoded := []byte(base64.StdEncoding.EncodeToString(b.Body))
	// Split into lines, max 76 bytes + crlf
	var buf bytes.Buffer
	for len(encoded) > 76 {
		buf.Write(encoded[:76])
		buf.Write(crlf)
		encoded = encoded[76:]
	}
	buf.Write(encoded)
	buf.Write(crlf)
	return buf.Bytes()
}

func (b MIMEBinary) Headers() map[string]string {
	return b.AdditionalHeaders
}

// Abstract MIME container for multipart bodies
type mIMEContainer struct {
	AdditionalHeaders map[string]string
	Parts             []MIMEPart
	boundary          []byte
}

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func (m *mIMEContainer) InitBoundary() {
	m.boundary = make([]byte, 32, 32)
	m.boundary[0] = byte('G')
	m.boundary[1] = byte('O')
	m.boundary[2] = byte('-')
	for i := 3; i < 32; i++ {
		m.boundary[i] = hexCharacters[randomGenerator.Intn(16)]
	}
}

func (m *mIMEContainer) AddPart(p MIMEPart) {
	m.Parts = append(m.Parts, p)
}

func (m mIMEContainer) Headers() map[string]string {
	return m.AdditionalHeaders
}

func (m mIMEContainer) ContentTransferEncoding() string {
	return ""
}

func (m mIMEContainer) EncodedBody() []byte {
	var buf bytes.Buffer
	buf.WriteString("This is a message with multiple parts in MIME format.")
	buf.Write(crlf)
	for _, part := range m.Parts {
		buf.WriteString("--")
		buf.Write(m.boundary)
		buf.Write(crlf)
		Write(part, &buf)
	}
	buf.WriteString("--")
	buf.Write(m.boundary)
	buf.WriteString("--")
	buf.Write(crlf)
	return buf.Bytes()
}

// Multipart/mixed is used for sending files with different "Content-Type" headers inline, or as attachments.
type MIMEMixed struct {
	mIMEContainer
}

func (m MIMEMixed) ContentType() string {
	if len(m.boundary) == 0 {
		log.Fatal("No MIME boundary set")
	}
	return "multipart/mixed; boundary=" + string(m.boundary)
}

// The multipart/alternative subtype indicates that each part is an "alternative" version of the same (or similar) content.
// The formats are ordered by how faithful they are to the original, with the least faithful first and the most faithful last.
// Most commonly, multipart/alternative is used for email with two parts, one plain text (text/plain) and one HTML (text/html).
type MIMEAlternative struct {
	mIMEContainer
}

func (m MIMEAlternative) ContentType() string {
	if len(m.boundary) == 0 {
		log.Fatal("No MIME boundary set")
	}
	return "multipart/alternative; boundary=" + string(m.boundary)
}

// Each message part is a component of an aggregate whole.
// The message consists of a root part which reference other parts inline.
// Message parts are commonly referenced by the "Content-ID" part header.
// One common usage of this subtype is to send a web page complete with images,
// in a single message. The root part would contain the HTML document,
// and use image tags to reference images stored in the latter parts.
type MIMERelated struct {
	mIMEContainer
}

func (m MIMERelated) ContentType() string {
	if len(m.boundary) == 0 {
		log.Fatal("No MIME boundary set")
	}
	return "multipart/related; boundary=" + string(m.boundary)
}
