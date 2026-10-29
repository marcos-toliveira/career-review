package cvdocx

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"
	"time"
)

// Options controls the document layout.
type Options struct {
	Lang  string // spell-check language, e.g. "pt-BR" or "en-US"
	Paper string // "a4" or "letter"
	Font  string // body font
	Title string // document title in the file properties
}

// Defaults for a given language: A4 for Portuguese, Letter for American English.
func Defaults(lang string) Options {
	o := Options{Lang: lang, Paper: "a4", Font: "Calibri"}
	if strings.EqualFold(lang, "en-US") {
		o.Paper = "letter"
	}
	return o
}

// fixedTime keeps the zip output byte-for-byte reproducible.
var fixedTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Write renders the blocks as a .docx package.
func Write(w io.Writer, blocks []Block, o Options) error {
	if o.Lang == "" {
		o.Lang = "pt-BR"
	}
	if o.Font == "" {
		o.Font = "Calibri"
	}
	width, height, err := paperSize(o.Paper)
	if err != nil {
		return err
	}
	parts := []struct{ name, body string }{
		{"[Content_Types].xml", contentTypes},
		{"_rels/.rels", rootRels},
		{"docProps/core.xml", coreProps(o)},
		{"word/_rels/document.xml.rels", documentRels},
		{"word/styles.xml", styles(o)},
		{"word/numbering.xml", numbering},
		{"word/document.xml", document(blocks, width, height)},
	}
	z := zip.NewWriter(w)
	for _, p := range parts {
		f, err := z.CreateHeader(&zip.FileHeader{Name: p.name, Method: zip.Deflate, Modified: fixedTime})
		if err != nil {
			return err
		}
		if _, err := io.WriteString(f, p.body); err != nil {
			return err
		}
	}
	return z.Close()
}

// paperSize returns the page size in twentieths of a point.
func paperSize(paper string) (int, int, error) {
	switch strings.ToLower(paper) {
	case "", "a4":
		return 11906, 16838, nil
	case "letter":
		return 12240, 15840, nil
	}
	return 0, 0, fmt.Errorf("papel desconhecido %q: use a4 ou letter", paper)
}

func document(blocks []Block, width, height int) string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`)
	if len(blocks) == 0 {
		b.WriteString("<w:p/>") // Word expects at least one paragraph in the body
	}
	for _, blk := range blocks {
		b.WriteString("<w:p><w:pPr>")
		switch blk.Kind {
		case Title:
			b.WriteString(`<w:pStyle w:val="Title"/>`)
		case Heading1:
			b.WriteString(`<w:pStyle w:val="Heading1"/>`)
		case Heading2:
			b.WriteString(`<w:pStyle w:val="Heading2"/>`)
		case Bullet:
			b.WriteString(`<w:pStyle w:val="ListBullet"/><w:numPr><w:ilvl w:val="0"/><w:numId w:val="1"/></w:numPr>`)
		}
		b.WriteString("</w:pPr>")
		for i, line := range blk.Lines {
			if i > 0 {
				b.WriteString("<w:r><w:br/></w:r>")
			}
			for _, s := range line {
				b.WriteString("<w:r>")
				if s.Bold || s.Italic {
					b.WriteString("<w:rPr>")
					if s.Bold {
						b.WriteString("<w:b/>")
					}
					if s.Italic {
						b.WriteString("<w:i/>")
					}
					b.WriteString("</w:rPr>")
				}
				fmt.Fprintf(&b, `<w:t xml:space="preserve">%s</w:t></w:r>`, escape(s.Text))
			}
		}
		b.WriteString("</w:p>")
	}
	// Margins: 2 cm on the sides, 1.8 cm top and bottom.
	fmt.Fprintf(&b, `<w:sectPr><w:pgSz w:w="%d" w:h="%d"/><w:pgMar w:top="1020" w:right="1134" w:bottom="1020" w:left="1134" w:header="0" w:footer="0" w:gutter="0"/></w:sectPr>`, width, height)
	b.WriteString("</w:body></w:document>")
	return b.String()
}

func styles(o Options) string {
	font := escape(o.Font)
	lang := escape(o.Lang)
	return xmlHeader + `<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:docDefaults><w:rPrDefault><w:rPr><w:rFonts w:ascii="` + font + `" w:hAnsi="` + font + `" w:cs="` + font + `" w:eastAsia="` + font + `"/>` +
		`<w:sz w:val="21"/><w:szCs w:val="21"/><w:lang w:val="` + lang + `"/></w:rPr></w:rPrDefault>` +
		`<w:pPrDefault><w:pPr><w:spacing w:after="80" w:line="264" w:lineRule="auto"/></w:pPr></w:pPrDefault></w:docDefaults>` +
		`<w:style w:type="paragraph" w:default="1" w:styleId="Normal"><w:name w:val="Normal"/><w:qFormat/></w:style>` +
		`<w:style w:type="paragraph" w:styleId="Title"><w:name w:val="Title"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
		`<w:pPr><w:spacing w:after="40"/></w:pPr><w:rPr><w:b/><w:color w:val="1F2937"/><w:sz w:val="40"/><w:szCs w:val="40"/></w:rPr></w:style>` +
		`<w:style w:type="paragraph" w:styleId="Heading1"><w:name w:val="heading 1"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
		`<w:pPr><w:keepNext/><w:spacing w:before="240" w:after="80"/><w:pBdr><w:bottom w:val="single" w:sz="6" w:space="2" w:color="9CA3AF"/></w:pBdr><w:outlineLvl w:val="0"/></w:pPr>` +
		`<w:rPr><w:b/><w:caps/><w:color w:val="1F2937"/><w:sz w:val="24"/><w:szCs w:val="24"/></w:rPr></w:style>` +
		`<w:style w:type="paragraph" w:styleId="Heading2"><w:name w:val="heading 2"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
		`<w:pPr><w:keepNext/><w:spacing w:before="160" w:after="40"/><w:outlineLvl w:val="1"/></w:pPr>` +
		`<w:rPr><w:b/><w:color w:val="111827"/><w:sz w:val="22"/><w:szCs w:val="22"/></w:rPr></w:style>` +
		`<w:style w:type="paragraph" w:styleId="ListBullet"><w:name w:val="List Bullet"/><w:basedOn w:val="Normal"/>` +
		`<w:pPr><w:spacing w:after="40"/><w:ind w:left="360" w:hanging="240"/></w:pPr></w:style>` +
		`</w:styles>`
}

func coreProps(o Options) string {
	return xmlHeader + `<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/">` +
		`<dc:title>` + escape(o.Title) + `</dc:title><dc:language>` + escape(o.Lang) + `</dc:language></cp:coreProperties>`
}

// escape writes s as XML text. Characters that XML 1.0 forbids (control characters such as
// form feeds pasted from PDFs) become spaces, and invalid UTF-8 becomes U+FFFD.
func escape(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == '&':
			b.WriteString("&amp;")
		case r == '<':
			b.WriteString("&lt;")
		case r == '>':
			b.WriteString("&gt;")
		case r == '"':
			b.WriteString("&quot;")
		case validXML(r):
			b.WriteRune(r)
		default:
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// validXML reports whether r is allowed in an XML 1.0 document.
func validXML(r rune) bool {
	return r == '\t' || r == '\n' || r == '\r' ||
		r >= 0x20 && r <= 0xD7FF || r >= 0xE000 && r <= 0xFFFD || r >= 0x10000 && r <= 0x10FFFF
}

const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

const contentTypes = xmlHeader + `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
	`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
	`<Default Extension="xml" ContentType="application/xml"/>` +
	`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
	`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>` +
	`<Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>` +
	`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>` +
	`</Types>`

const rootRels = xmlHeader + `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
	`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
	`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>` +
	`</Relationships>`

const documentRels = xmlHeader + `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
	`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>` +
	`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"/>` +
	`</Relationships>`

const numbering = xmlHeader + `<w:numbering xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
	`<w:abstractNum w:abstractNumId="0"><w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="bullet"/><w:lvlText w:val="•"/><w:lvlJc w:val="left"/>` +
	`<w:pPr><w:ind w:left="360" w:hanging="240"/></w:pPr></w:lvl></w:abstractNum>` +
	`<w:num w:numId="1"><w:abstractNumId w:val="0"/></w:num></w:numbering>`
