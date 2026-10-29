package cvdocx

import (
	"reflect"
	"testing"
)

func TestParseBlocks(t *testing.T) {
	src := "# Nome\n\nCargo · área\nCidade · e-mail\n\n## Experiência\n\n### Empresa · 2020\n\n- item um\n  continua\n- item dois\n"
	got := Parse(src)
	want := []Block{
		{Kind: Title, Lines: []Line{{{Text: "Nome"}}}},
		{Kind: Paragraph, Lines: []Line{{{Text: "Cargo · área"}}, {{Text: "Cidade · e-mail"}}}},
		{Kind: Heading1, Lines: []Line{{{Text: "Experiência"}}}},
		{Kind: Heading2, Lines: []Line{{{Text: "Empresa · 2020"}}}},
		{Kind: Bullet, Lines: []Line{{{Text: "item um"}, {Text: " "}, {Text: "continua"}}}},
		{Kind: Bullet, Lines: []Line{{{Text: "item dois"}}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse:\n got %#v\nwant %#v", got, want)
	}
}

func TestInline(t *testing.T) {
	cases := []struct {
		in   string
		want Line
	}{
		{"**Cargo** e texto", Line{{Text: "Cargo", Bold: true}, {Text: " e texto"}}},
		{"um *itálico* aqui", Line{{Text: "um "}, {Text: "itálico", Italic: true}, {Text: " aqui"}}},
		{"snake_case_name", Line{{Text: "snake_case_name"}}},
		{"2 * 3", Line{{Text: "2 * 3"}}},
		{"[meu site](https://exemplo.com) e `code`", Line{{Text: "meu site e code"}}},
		{"colchete [solto", Line{{Text: "colchete [solto"}}},
	}
	for _, c := range cases {
		if got := inline(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("inline(%q):\n got %#v\nwant %#v", c.in, got, c.want)
		}
	}
}
