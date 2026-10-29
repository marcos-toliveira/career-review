package cvdocx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

func render(t *testing.T, src string, o Options) map[string]string {
	t.Helper()
	var buf bytes.Buffer
	if err := Write(&buf, Parse(src), o); err != nil {
		t.Fatal(err)
	}
	z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	parts := map[string]string{}
	for _, f := range z.File {
		r, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(r)
		r.Close()
		parts[f.Name] = string(b)
	}
	return parts
}

func TestWritePackageIsWellFormed(t *testing.T) {
	parts := render(t, "# A & B <C>\n\n- **x** \"y\"\n", Defaults("pt-BR"))
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "docProps/core.xml", "word/_rels/document.xml.rels", "word/styles.xml", "word/numbering.xml", "word/document.xml"} {
		body, ok := parts[name]
		if !ok {
			t.Fatalf("parte ausente: %s", name)
		}
		d := xml.NewDecoder(strings.NewReader(body))
		for {
			if _, err := d.Token(); err == io.EOF {
				break
			} else if err != nil {
				t.Fatalf("%s não é XML válido: %v", name, err)
			}
		}
	}
	doc := parts["word/document.xml"]
	for _, want := range []string{"A &amp; B &lt;C&gt;", `<w:pStyle w:val="Title"/>`, `<w:numId w:val="1"/>`, "<w:b/>", "&quot;y&quot;"} {
		if !strings.Contains(doc, want) {
			t.Errorf("document.xml sem %q", want)
		}
	}
}

func TestLanguageAndPaper(t *testing.T) {
	pt := render(t, "texto", Defaults("pt-BR"))
	en := render(t, "text", Defaults("en-US"))
	if !strings.Contains(pt["word/styles.xml"], `w:lang w:val="pt-BR"`) || !strings.Contains(en["word/styles.xml"], `w:lang w:val="en-US"`) {
		t.Fatal("idioma do corretor não aplicado")
	}
	if !strings.Contains(pt["word/document.xml"], `w:w="11906"`) || !strings.Contains(en["word/document.xml"], `w:w="12240"`) {
		t.Fatal("papel padrão errado: esperado A4 em pt-BR e Letter em en-US")
	}
	var buf bytes.Buffer
	if err := Write(&buf, nil, Options{Paper: "ofício"}); err == nil {
		t.Fatal("papel desconhecido deveria falhar")
	}
}

func TestDeterministic(t *testing.T) {
	var a, b bytes.Buffer
	blocks := Parse("# Nome\n\n- item\n")
	if err := Write(&a, blocks, Defaults("pt-BR")); err != nil {
		t.Fatal(err)
	}
	if err := Write(&b, blocks, Defaults("pt-BR")); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Fatal("a mesma entrada gerou arquivos diferentes")
	}
}
