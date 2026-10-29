package cvdocx

import (
	"reflect"
	"testing"
)

func TestParseBlocks(t *testing.T) {
	src := "# Nome\n\nCargo · área\r\nCidade · e-mail\n\n## Experiência\n\n### Empresa · 2020\n\n- item um\n  continua\n- item dois\n\tcom tab\n\n1. primeiro\n2. segundo\n\n#### Fundo\n\n#hashtag\n"
	got := Parse(src)
	want := []Block{
		{Kind: Title, Lines: []Line{{{Text: "Nome"}}}},
		{Kind: Paragraph, Lines: []Line{{{Text: "Cargo · área"}}, {{Text: "Cidade · e-mail"}}}},
		{Kind: Heading1, Lines: []Line{{{Text: "Experiência"}}}},
		{Kind: Heading2, Lines: []Line{{{Text: "Empresa · 2020"}}}},
		{Kind: Bullet, Lines: []Line{{{Text: "item um"}, {Text: " "}, {Text: "continua"}}}},
		{Kind: Bullet, Lines: []Line{{{Text: "item dois"}, {Text: " "}, {Text: "com tab"}}}},
		{Kind: Paragraph, Lines: []Line{{{Text: "1. primeiro"}}}},
		{Kind: Paragraph, Lines: []Line{{{Text: "2. segundo"}}}},
		{Kind: Heading2, Lines: []Line{{{Text: "Fundo"}}}},
		{Kind: Paragraph, Lines: []Line{{{Text: "#hashtag"}}}},
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
		{"__Cargo__: x", Line{{Text: "Cargo", Bold: true}, {Text: ": x"}}},
		{"um *itálico* aqui", Line{{Text: "um "}, {Text: "itálico", Italic: true}, {Text: " aqui"}}},
		{"_Cargo_—atual", Line{{Text: "Cargo", Italic: true}, {Text: "—atual"}}},
		{"**Ênfase**, depois", Line{{Text: "Ênfase", Bold: true}, {Text: ", depois"}}},
		{"snake_case_name", Line{{Text: "snake_case_name"}}},
		{"2 * 3", Line{{Text: "2 * 3"}}},
		{"custo 2*3 vezes", Line{{Text: "custo 2*3 vezes"}}},
		{"salário* a confirmar", Line{{Text: "salário* a confirmar"}}},
		{"Texto **sem fechamento", Line{{Text: "Texto **sem fechamento"}}},
		{"*a* e *b*", Line{{Text: "a", Italic: true}, {Text: " e "}, {Text: "b", Italic: true}}},
		{"[meu site](https://exemplo.com) e `code`", Line{{Text: "meu site (exemplo.com) e code"}}},
		{"[exemplo.com](https://exemplo.com)", Line{{Text: "exemplo.com"}}},
		{"[Certificação] Ver [Credencial](https://x.io)", Line{{Text: "[Certificação] Ver Credencial (x.io)"}}},
		{"colchete [solto", Line{{Text: "colchete [solto"}}},
		{"[a [b](https://u.io)", Line{{Text: "[a b (u.io)"}}},
		{"crase `solta", Line{{Text: "crase `solta"}}},
	}
	for _, c := range cases {
		if got := inline(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("inline(%q):\n got %#v\nwant %#v", c.in, got, c.want)
		}
	}
}
