package plugin

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"

	"golang.org/x/tools/go/loader"
	"gopkg.in/yaml.v2"
)

func cleanPlugins() error {
	os.Remove(os.Getenv("GOPATH") + "/src/github.com/gizak/gopp/plugin/plugins.go")
	return nil
}

func installPlugins(list PList) error {

	for i, p := range list {
		if p.Plugin == "" {
			return errors.New("Bad entry at " + strconv.Itoa(i))
		}

		cmd := exec.Command("go", "get", p.Uri)
		log.Println("LOAD plugin: " + p.Plugin)
		if b, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprint(os.Stderr, string(b))
			return err
		}
	}

	return nil
}

func checkPlugins(list PList) error {

	var conf loader.Config
	pluginExport := "PLUGIN"

	// chk exported PLUGIN
	for _, p := range list {
		// check parsing
		fs := token.NewFileSet()
		pkgs, err := parser.ParseDir(fs, os.Getenv("GOPATH")+"/src/"+p.Uri, nil, 0)
		if err != nil {
			return err
		}
		if len(pkgs) != 1 {
			return errors.New("Bad package definition: " + p.Plugin)
		}

		var pkg *ast.Package
		for pname := range pkgs {
			pkg = pkgs[pname]
			log.Println("CHECK plugin " + p.Plugin + ", pkg:" + pname)
		}
		ast.PackageExports(pkg)

		// check exported symbol
		found := false
		ast.Inspect(pkg, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.ValueSpec:
				if x.Names[0].Name == pluginExport {
					found = true
					return false
				}
			}
			return true
		})
		if !found {
			return errors.New("No exported plugin found for " + p.Plugin + " (expected declare: var PLUGIN)")
		}

		// check others
		conf.Import(p.Uri)
	}

	_, err := conf.Load()
	return err
}

func buildPlugins(list PList) error {

	gopath := path.Clean(os.Getenv("GOPATH"))
	genfile := path.Join(gopath, "src/github.com/gizak/gopp/plugin/plugins.go")

	pre := "package plugin\n"

	imps := ""
	init := "func init() {\n "
	for i, p := range list {

		log.Println("WRITE plugin entry: " + p.Plugin)

		pname := "p" + strconv.Itoa(i)
		imps += fmt.Sprintf("import %s %q\n", pname, p.Uri)
		init += fmt.Sprintf("\tPlugins = append(Plugins,PluginEntry{Name:%q,Plugin:%s.PLUGIN,Uri:%q})\n", p.Plugin, pname, p.Uri)
	}
	init += "}"

	if err := ioutil.WriteFile(genfile, []byte(pre+imps+init), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func Init() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	log.SetPrefix("> [init]\t")
	regPlugins()
}

func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}

func Install(list PList) error {
	// clean installed plugins
	if err := cleanPlugins(); err != nil {
		return err
	}

	// install plugins
	if err := installPlugins(list); err != nil {
		return err
	}

	// parsing check
	if err := checkPlugins(list); err != nil {
		return err
	}

	// build plugins from source
	if err := buildPlugins(list); err != nil {
		return err
	}

	regPlugins()

	log.Println("BUILD gopp")
	cmd := exec.Command("go", "install", "github.com/gizak/gopp")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + string(out))
	}

	return nil
}

func InstallWithConfig(fpath string) error {
	log.SetPrefix("➧ [install]\t")

	log.Println("READ config: " + fpath)
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return Install(cfg.Plugins)
}
