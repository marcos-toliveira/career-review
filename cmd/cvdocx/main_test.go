package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGeneratesDocxNextToInput(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "cv.md")
	if err := os.WriteFile(in, []byte("# Ana\n\n- item\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-lang", "en-US", in}); err != nil {
		t.Fatal(err)
	}
	z, err := zip.OpenReader(filepath.Join(dir, "cv.docx"))
	if err != nil {
		t.Fatalf("saída não é um .docx: %v", err)
	}
	z.Close()
}

func TestRunErrors(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "cv.md")
	os.WriteFile(in, []byte("# Ana\n"), 0o644)
	cases := map[string][]string{
		"sem arquivo":         {},
		"arquivo ausente":     {filepath.Join(dir, "nao-existe.md")},
		"idioma inválido":     {"-lang", "portugues", in},
		"papel inválido":      {"-paper", "oficio", in},
		"saída igual à fonte": {"-o", in, in},
	}
	for name, args := range cases {
		if err := run(args); err == nil {
			t.Errorf("%s: esperava erro", name)
		}
	}
	if b, _ := os.ReadFile(in); !strings.HasPrefix(string(b), "# Ana") {
		t.Fatal("o arquivo de entrada foi alterado")
	}
}

func TestExamples(t *testing.T) {
	for _, ex := range []struct{ file, lang string }{
		{"curriculo-exemplo-ptbr.md", "pt-BR"},
		{"curriculo-exemplo-en.md", "en-US"},
	} {
		out := filepath.Join(t.TempDir(), "ex.docx")
		if err := run([]string{"-lang", ex.lang, "-o", out, filepath.Join("..", "..", "examples", ex.file)}); err != nil {
			t.Fatalf("%s: %v", ex.file, err)
		}
	}
}
