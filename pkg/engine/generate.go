package engine

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LAShZ/go-scaffold/config"
	"github.com/LAShZ/go-scaffold/pkg/tempfs"
)

var ProjectName string
var ProjectModule string
var GoVersion string
var Data map[string]interface{}
var Verbose bool

type ProjectGenerator struct {
	fs     tempfs.TempFS
	logger string
	orm    string
	web    string
	redis  bool
}

type Options func(*ProjectGenerator)

func NewProjectGenerator(fs tempfs.TempFS, opts ...Options) *ProjectGenerator {
	pg := &ProjectGenerator{
		fs: fs,
	}
	for _, o := range opts {
		o(pg)
	}
	return pg
}

func WithLogger(logger string) Options {
	return func(pg *ProjectGenerator) {
		pg.logger = logger
	}
}

func WithORM(orm string) Options {
	return func(pg *ProjectGenerator) {
		pg.orm = orm
	}
}

func WithWeb(web string) Options {
	return func(pg *ProjectGenerator) {
		pg.web = web
	}
}

func UseRedis(redis bool) Options {
	return func(pg *ProjectGenerator) {
		pg.redis = redis
	}
}

func (pg *ProjectGenerator) Generate() {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("Generate project failed, err: %s\n", err)
			return
		}
		gofmt := exec.Command("go", "fmt", "./...")
		err = gofmt.Run()
		if err != nil {
			fmt.Println("fmt project failed,", err)
		}
		fmt.Printf("\n\n\033[0;32mProject created: \033[0m%s \n\n", ProjectName)
		fmt.Printf("Run:\tmake dep\n\tmake build\n\tto build the project\n\n")
		fmt.Printf("Use: ./%s to start the project\n", ProjectName)
	}()

	Data = make(map[string]interface{})
	Data["ProjectModule"] = ProjectModule
	Data["ProjectName"] = ProjectName
	Data["GoVersion"] = GoVersion
	Data["Config"] = config.Info

	var entrys []fs.DirEntry
	entrys, err = pg.fs.ReadDir("template")
	for _, entry := range entrys {
		if entry.IsDir() {
			err = pg.generateDir(entry, "template")
			if err != nil {
				return
			}
		} else {
			tmplname := entry.Name()
			filename := strings.TrimSuffix(tmplname, ".tmpl")
			err = pg.generateFile(filename, "template"+string(filepath.Separator)+tmplname, Data)
			if err != nil {
				return
			}
		}
	}
}

func (pg *ProjectGenerator) generateFile(filename string, tmplname string, data map[string]interface{}) error {
	// since go embed not support embed file has prefix '.' or '_',
	// file has name like that should trim '.' or '_' in template folder,
	// and add the prefix here
	if filename == "gitignore" || filename == "dockerignore" || filename == "golangci.yml" {
		filename = "." + filename
	}

	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	tmpl, err := pg.fs.Open(tmplname)
	if err != nil {
		return err
	}
	text, err := ioutil.ReadAll(tmpl)
	if err != nil {
		return err
	}
	if Verbose {
		fmt.Printf("\n\n\033[0;32mApplying tmpl: \033[0m%s \n %s\n\n", tmplname, text)
	}
	fileTmpl, err := template.New(filename).Parse(string(text))
	if err != nil {
		return err
	}
	err = fileTmpl.Execute(fd, data)
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(fileData) == 0 {
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg *ProjectGenerator) generateDir(entry fs.DirEntry, prefix string) error {
	err := os.Mkdir(entry.Name(), 0755)
	if err != nil {
		return err
	}

	err = os.Chdir(entry.Name())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			fmt.Printf("Generate project failed, err: %s\n", err)
		}
		err := os.Chdir("../")
		if err != nil {
			fmt.Printf("Generate project failed, err: %s\n", err)
		}
	}()

	path := prefix + string(filepath.Separator) + entry.Name()
	entrys, err := pg.fs.ReadDir(path)

	if err != nil {
		return err
	}

	for _, entry := range entrys {
		if entry.IsDir() {
			err = pg.generateDir(entry, path)
			if err != nil {
				return err
			}
			tempEntrys, err := os.ReadDir(entry.Name())
			if err != nil {
				return err
			}
			if len(tempEntrys) == 0 {
				_ = os.RemoveAll(entry.Name())
			}
			continue
		}
		tmplname := entry.Name()
		filename := strings.TrimSuffix(tmplname, ".tmpl")
		err = pg.generateFile(filename, path+string(filepath.Separator)+tmplname, Data)
		if err != nil {
			return err
		}
	}

	return nil
}
