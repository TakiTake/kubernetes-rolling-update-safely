package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Variables struct{}

func getWriter(out string) io.Writer {
	if out == "" {
		return os.Stdout
	}

	w, err := os.Create(out)
	if err != nil {
		panic(err)
	}

	return w
}

func include(f string) string {
	d, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return string(d)
}

func diff(f1, f2 string) string {
	out, _ := exec.Command("diff", f1, f2).CombinedOutput()
	return string(out)
}

func safeHTML(val string) template.HTML {
	return template.HTML(val)
}

func makeReadme(w io.Writer, v *Variables) {
	funcMap := template.FuncMap{
		"include":  include,
		"diff":     diff,
		"safeHTML": safeHTML,
	}

	tmpl := template.Must(template.New("README.md.tmpl").Funcs(funcMap).ParseFiles("README.md.tmpl"))
	err := tmpl.Execute(w, v)
	if err != nil {
		panic(err)
	}
}

func main() {
	var out string
	var v = &Variables{}

	flag.StringVar(&out, "o", "", "Write output to <file> instead of stdout")
	flag.Parse()

	w := getWriter(out)

	makeReadme(w, v)
}
