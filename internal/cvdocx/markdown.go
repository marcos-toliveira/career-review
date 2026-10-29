// Package cvdocx converts the small Markdown subset used by résumés into a .docx file.
package cvdocx

import "strings"

// Kind identifies the type of a block.
type Kind int

const (
	Title Kind = iota
	Heading1
	Heading2
	Paragraph
	Bullet
)

// Span is a run of text with optional bold or italic formatting.
type Span struct {
	Text   string
	Bold   bool
	Italic bool
}

// Line is one visual line of a block; lines of the same block are separated by line breaks.
type Line []Span

// Block is a paragraph-level element.
type Block struct {
	Kind  Kind
	Lines []Line
}

// Parse reads Markdown and returns its blocks. Supported: "#", "##" and "###" headings,
// "-" or "*" bullets, paragraphs, **bold**, *italic* or _italic_, `code` and [text](url) links.
// Lines inside a paragraph keep their line breaks, which suits résumé layouts.
func Parse(src string) []Block {
	var blocks []Block
	var cur *Block
	flush := func() {
		if cur != nil {
			blocks = append(blocks, *cur)
			cur = nil
		}
	}
	for _, raw := range strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n") {
		line := strings.TrimRight(raw, " \t")
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			flush()
		case strings.HasPrefix(trimmed, "### "):
			flush()
			blocks = append(blocks, Block{Kind: Heading2, Lines: []Line{inline(trimmed[4:])}})
		case strings.HasPrefix(trimmed, "## "):
			flush()
			blocks = append(blocks, Block{Kind: Heading1, Lines: []Line{inline(trimmed[3:])}})
		case strings.HasPrefix(trimmed, "# "):
			flush()
			blocks = append(blocks, Block{Kind: Title, Lines: []Line{inline(trimmed[2:])}})
		case strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* "):
			flush()
			cur = &Block{Kind: Bullet, Lines: []Line{inline(trimmed[2:])}}
		case cur != nil && cur.Kind == Bullet && strings.HasPrefix(line, " "):
			last := &cur.Lines[len(cur.Lines)-1]
			*last = append(*last, Span{Text: " "})
			*last = append(*last, inline(trimmed)...)
		default:
			if cur == nil || cur.Kind != Paragraph {
				flush()
				cur = &Block{Kind: Paragraph}
			}
			cur.Lines = append(cur.Lines, inline(trimmed))
		}
	}
	flush()
	return blocks
}

// inline splits text into spans, handling bold, italic, code and links.
func inline(s string) Line {
	var out Line
	var buf strings.Builder
	bold, italic := false, false
	emit := func() {
		if buf.Len() > 0 {
			out = append(out, Span{Text: buf.String(), Bold: bold, Italic: italic})
			buf.Reset()
		}
	}
	for i := 0; i < len(s); {
		switch {
		case strings.HasPrefix(s[i:], "**") || strings.HasPrefix(s[i:], "__"):
			emit()
			bold = !bold
			i += 2
		case (s[i] == '*' || s[i] == '_') && italicDelimiter(s, i):
			emit()
			italic = !italic
			i++
		case s[i] == '`':
			if end := strings.IndexByte(s[i+1:], '`'); end >= 0 {
				buf.WriteString(s[i+1 : i+1+end])
				i += end + 2
			} else {
				buf.WriteByte(s[i])
				i++
			}
		case s[i] == '[':
			text, n := link(s[i:])
			if n == 0 {
				buf.WriteByte(s[i])
				i++
				continue
			}
			buf.WriteString(text)
			i += n
		default:
			buf.WriteByte(s[i])
			i++
		}
	}
	emit()
	return out
}

// italicDelimiter reports whether the single '*' or '_' at i opens or closes emphasis,
// so that snake_case words and "2 * 3" stay literal.
func italicDelimiter(s string, i int) bool {
	prevWord := i > 0 && isWordByte(s[i-1])
	nextWord := i+1 < len(s) && isWordByte(s[i+1])
	if s[i] == '_' {
		return prevWord != nextWord
	}
	nextSpace := i+1 >= len(s) || s[i+1] == ' '
	prevSpace := i == 0 || s[i-1] == ' '
	return !(nextSpace && prevSpace)
}

func isWordByte(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b >= 0x80
}

// link parses "[text](url)" at the start of s and returns the visible text and the bytes consumed.
func link(s string) (string, int) {
	closeText := strings.Index(s, "](")
	if closeText < 0 {
		return "", 0
	}
	closeURL := strings.IndexByte(s[closeText+2:], ')')
	if closeURL < 0 {
		return "", 0
	}
	return s[1:closeText], closeText + 2 + closeURL + 1
}
