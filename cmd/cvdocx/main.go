// Command cvdocx converts a résumé written in Markdown into a .docx file.
//
//	cvdocx -lang pt-BR curriculo-ptbr.md            # gera curriculo-ptbr.docx
//	cvdocx -lang en-US -o resume.docx resume-en.md
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/marcos-toliveira/career-review/internal/cvdocx"
)

var langTag = regexp.MustCompile(`^[a-z]{2,3}(-[A-Z]{2})?$`)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "cvdocx:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("cvdocx", flag.ContinueOnError)
	lang := fs.String("lang", "pt-BR", "idioma do corretor ortográfico: pt-BR ou en-US")
	paper := fs.String("paper", "", "papel: a4 ou letter (padrão: a4 para pt-BR, letter para en-US)")
	font := fs.String("font", "Calibri", "fonte do texto")
	out := fs.String("o", "", "arquivo de saída (padrão: o nome do .md com extensão .docx)")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "uso: cvdocx [-lang pt-BR|en-US] [-paper a4|letter] [-font nome] [-o saida.docx] curriculo.md")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("informe um arquivo .md")
	}
	in := fs.Arg(0)
	if !langTag.MatchString(*lang) {
		return fmt.Errorf("idioma %q inválido: use um código como pt-BR ou en-US", *lang)
	}
	src, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	o := cvdocx.Defaults(*lang)
	if *paper != "" {
		o.Paper = *paper
	}
	o.Font = *font
	blocks := cvdocx.Parse(string(src))
	for _, b := range blocks {
		if b.Kind == cvdocx.Title && len(b.Lines) > 0 {
			for _, s := range b.Lines[0] {
				o.Title += s.Text
			}
			break
		}
	}
	dst := *out
	if dst == "" {
		dst = strings.TrimSuffix(in, ".md") + ".docx"
	}
	if filepath.Clean(dst) == filepath.Clean(in) {
		return fmt.Errorf("a saída %q é o próprio arquivo de entrada", dst)
	}
	var buf bytes.Buffer
	if err := cvdocx.Write(&buf, blocks, o); err != nil {
		return err
	}
	if err := os.WriteFile(dst, buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Println(dst)
	return nil
}
