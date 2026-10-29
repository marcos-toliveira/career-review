// Package cvdocx converts the small Markdown subset used by résumés into a .docx file.
package cvdocx

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

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

var numbered = regexp.MustCompile(`^\d+[.)] `)

// Parse reads Markdown and returns its blocks. Supported: "#" (title), "##" and "###" headings
// (deeper levels are treated as "###"), "-" or "*" bullets, numbered items (kept as their own
// paragraph, with the number), paragraphs, **bold**, *italic* or _italic_, `code` and
// [text](url) links, written as "text (url)". Lines inside a paragraph keep their line breaks,
// which suits résumé layouts.
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
		line := strings.TrimRight(raw, " \t\r")
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			flush()
		case strings.HasPrefix(trimmed, "#"):
			level := len(trimmed) - len(strings.TrimLeft(trimmed, "#"))
			text := strings.TrimSpace(trimmed[level:])
			if text == "" || trimmed[level] != ' ' {
				goto paragraph
			}
			flush()
			kind := Heading2
			switch level {
			case 1:
				kind = Title
			case 2:
				kind = Heading1
			}
			blocks = append(blocks, Block{Kind: kind, Lines: []Line{inline(text)}})
		case strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* "):
			flush()
			cur = &Block{Kind: Bullet, Lines: []Line{inline(trimmed[2:])}}
		case numbered.MatchString(trimmed):
			flush()
			blocks = append(blocks, Block{Kind: Paragraph, Lines: []Line{inline(trimmed)}})
		case cur != nil && cur.Kind == Bullet && (line[0] == ' ' || line[0] == '\t'):
			last := &cur.Lines[len(cur.Lines)-1]
			*last = append(*last, Span{Text: " "})
			*last = append(*last, inline(trimmed)...)
		default:
			goto paragraph
		}
		continue
	paragraph:
		if cur == nil || cur.Kind != Paragraph {
			flush()
			cur = &Block{Kind: Paragraph}
		}
		cur.Lines = append(cur.Lines, inline(trimmed))
	}
	flush()
	return blocks
}

// inline splits text into spans, handling bold, italic, code and links. An emphasis
// delimiter only counts when it has a matching partner on the same line; otherwise it is
// kept as a literal character, so "2*3", "nota*" and "snake_case" survive unchanged.
func inline(s string) Line {
	var out Line
	var buf strings.Builder
	bold, italic := "", "" // the delimiter that opened each style, or "" when closed
	emit := func() {
		if buf.Len() > 0 {
			out = append(out, Span{Text: buf.String(), Bold: bold != "", Italic: italic != ""})
			buf.Reset()
		}
	}
	for i := 0; i < len(s); {
		switch {
		case s[i] == '*' || s[i] == '_':
			d := s[i : i+1]
			if strings.HasPrefix(s[i:], d+d) {
				d += d
			}
			open := &italic
			if len(d) == 2 {
				open = &bold
			}
			switch {
			case *open == d && closes(s, i, len(d)):
				emit()
				*open = ""
			case *open == "" && opens(s, i, len(d)) && hasCloser(s, i+len(d), d):
				emit()
				*open = d
			default:
				buf.WriteString(d)
			}
			i += len(d)
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

// opens reports whether the delimiter of length n at i can start emphasis:
// followed by a non-space and not preceded by a letter or digit.
func opens(s string, i, n int) bool {
	next, _ := utf8.DecodeRuneInString(s[i+n:])
	prev, _ := utf8.DecodeLastRuneInString(s[:i])
	return i+n < len(s) && !unicode.IsSpace(next) && (i == 0 || !isWord(prev))
}

// closes reports whether the delimiter of length n at i can end emphasis:
// preceded by a non-space and not followed by a letter or digit.
func closes(s string, i, n int) bool {
	prev, _ := utf8.DecodeLastRuneInString(s[:i])
	next, _ := utf8.DecodeRuneInString(s[i+n:])
	return i > 0 && !unicode.IsSpace(prev) && (i+n == len(s) || !isWord(next))
}

// hasCloser reports whether a closing delimiter d appears in s at or after from.
func hasCloser(s string, from int, d string) bool {
	for j := from; j < len(s); j++ {
		if !strings.HasPrefix(s[j:], d) {
			continue
		}
		// A single delimiter must not be half of a double one.
		if len(d) == 1 && (strings.HasPrefix(s[j+1:], d) || (j > 0 && s[j-1] == d[0])) {
			continue
		}
		if closes(s, j, len(d)) {
			return true
		}
	}
	return false
}

func isWord(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// link parses "[text](url)" at the start of s and returns the visible text and the bytes
// consumed. The URL is kept after the text unless the text already shows it.
func link(s string) (string, int) {
	closeText := strings.IndexByte(s, ']')
	if closeText < 0 || strings.IndexByte(s[1:closeText], '[') >= 0 || !strings.HasPrefix(s[closeText:], "](") {
		return "", 0
	}
	closeURL := matchingParen(s[closeText+2:])
	if closeURL < 0 {
		return "", 0
	}
	text, url := s[1:closeText], s[closeText+2:closeText+2+closeURL]
	n := closeText + 2 + closeURL + 1
	bare := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(url, "mailto:"), "https://"), "http://")
	if url == "" || text == url || text == bare {
		return text, n
	}
	return text + " (" + bare + ")", n
}

// matchingParen returns the index of the ")" that closes a link destination, allowing
// balanced parentheses inside it, or -1.
func matchingParen(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			if depth == 0 {
				return i
			}
			depth--
		}
	}
	return -1
}
