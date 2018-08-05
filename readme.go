package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Variables struct {
	V1_0           string
	Diff_V1_0_V1_1 string
}

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

func setVariables(v *Variables) {
	v1_0, err := ioutil.ReadFile("demo-manifest/1-0.default.deploy.yml")
	if err != nil {
		panic(err)
	}
	v.V1_0 = string(v1_0)

	diff, _ := exec.Command("diff", "demo-manifest/1-0.default.deploy.yml", "demo-manifest/1-1.default.deploy.yml").CombinedOutput()
	v.Diff_V1_0_V1_1 = string(diff)
}

func safeHTML(val string) template.HTML {
	return template.HTML(val)
}

func makeReadme(w io.Writer, v *Variables) {
	funcMap := template.FuncMap{
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
	setVariables(v)

	makeReadme(w, v)
}
