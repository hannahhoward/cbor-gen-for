package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

var ErrWrongPackageCount = errors.New("Must be part of only one package")

type Generator struct {
	Path       string
	Filename   string
	Package    string
	GenStructs []string
}

type templateData struct {
	Path       string
	Filebase   string
	Package    string
	GenStructs []string
	PkgPath    string
}

func (g Generator) GenerateCborTypes() error {
	fpath := filepath.Join(g.Path, g.Filename)
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName}, "file="+fpath)
	if err != nil {
		return err
	}

	if len(pkgs) != 1 {
		return ErrWrongPackageCount
	}
	tdata := templateData{
		Path:       g.Path,
		Package:    g.Package,
		Filebase:   strings.TrimSuffix(g.Filename, ".go"),
		PkgPath:    pkgs[0].PkgPath,
		GenStructs: g.GenStructs,
	}

	t, err := template.New("run_cbor_gen").Parse(defaultTemplate)
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "gen")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	f, err := os.OpenFile(filepath.Join(dir, "main.go"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	err = t.Execute(f, tdata)
	if err != nil {
		return err
	}
	f.Close()
	cmd := exec.Command("go", "run", filepath.Join(dir, "main.go"))
	output, err := cmd.Output()
	fmt.Println(string(output))
	return err
}

var (
	defaultTemplate = `
package main

import (
	"fmt"
	"os"

	cborgen "github.com/whyrusleeping/cbor-gen"
	{{.Package}} "{{.PkgPath}}"
)

func main() {
	fmt.Print("Generating Cbor Marshal/Unmarshal...")

	genName := "{{.Path}}/{{.Filebase}}_cbor_gen.go"
	if err := cborgen.WriteTupleEncodersToFile(
		genName,
		"{{.Package}}",
		{{$gen := .}}
		{{range .GenStructs}}
			{{$gen.Package}}.{{.}}{},
		{{end}}
	); err != nil {
		fmt.Println("Failed: ")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done.")
}
`
)
