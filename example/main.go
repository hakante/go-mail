package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	. "terelius.dev/go/mail"
)

func main() {
	from := Address{"Me", "me@example.com"}
	to := []Address{Address{"You", "you@example.com"}}
	subject := "Hello World"
	someText := `This is the message content.
\n

Technically speaking, the Unicode characters are embedded in 8 bit HTML using 'character entities', for instance:

&#2384; = ॐ
&#1488; = א‎
&#937; = Ω
If your browser is Unicode-enabled, you should see the Sanskrit letter for 'Aum' (see this image); the Hebrew letter Aleph, and a Greek capital Omega above.`
	someHTML := `<!doctype html>
<html lang=en>
<head>
<meta charset=utf-8>
<title>This is the title</title>
</head>
<body>
<p>I'm the content</p>
<img src="cid:Smiley.jpg">
<p><a href="www.google.com">Google</a></p>
</body>
</html>`

	fileBytes, err := ioutil.ReadFile("example/Smiley.jpg")
	if err != nil {
		log.Fatal(err)
	}
	imageAttachment := MIMEBinary{Type: mime.TypeByExtension(".jpg"), Body: fileBytes, AdditionalHeaders: map[string]string{"Content-Disposition": "inline; filename=Smiley.jpg", "Content-ID": "<Smiley.jpg>"}}

	related := MIMERelated{}
	related.InitBoundary()
	related.AddPart(MIMEHTML(someHTML))
	related.AddPart(imageAttachment)

	alternative := MIMEAlternative{}
	alternative.InitBoundary()
	alternative.AddPart(MIMEPlain(someText))
	alternative.AddPart(related)

	mixed := MIMEMixed{}
	mixed.InitBoundary()
	mixed.AddPart(alternative)
	mixed.AddPart(MIMEPlain("And here follows some more text after the smiley"))

	message := Message{From: from, To: to, Subject: subject, Content: mixed}
	fmt.Println(string(message.String()))
}
