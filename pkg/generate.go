package pkg

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/LAShZ/go-scaffold/config"
)

var ProjectName string
var ProjectModule string
var GoVersion string
var Data map[string]interface{}
var Verbose bool

type ProjectGenerator struct {
	fs     embed.FS
	logger string
	orm    string
	web    string
	redis  bool
}

type Options func(*ProjectGenerator)

func NewProjectGenerator(fs embed.FS, opts ...Options) *ProjectGenerator {
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
		}
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

func (pg *ProjectGenerator) generateFile(filename string, tmplname string, data map[string]interface{}) (err error) {
	fd, err := os.Create(filename)
	if err != nil {
		return
	}
	tmpl, err := pg.fs.Open(tmplname)
	if err != nil {
		return
	}
	text, err := ioutil.ReadAll(tmpl)
	if err != nil {
		return
	}
	if Verbose {
		fmt.Printf("\n\n\033[0;32mApplying tmpl: \033[0m%s \n %s\n\n", tmplname, text)
	}
	fileTmpl, err := template.New(filename).Parse(string(text))
	if err != nil {
		return
	}
	err = fileTmpl.Execute(fd, data)
	return
}

func (pg *ProjectGenerator) generateDir(entry fs.DirEntry, prefix string) error {
	err := os.Mkdir(entry.Name(), 0755)
	if err != nil {
		return err
	}

	err = os.Chdir(entry.Name())
	if err != nil {
		return err
	}
	defer func() {
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

	fmt.Println("entry:", entrys)

	for _, entry := range entrys {
		if entry.IsDir() {
			err = pg.generateDir(entry, path)
			if err != nil {
				return err
			}
		} else {
			tmplname := entry.Name()
			filename := strings.TrimSuffix(tmplname, ".tmpl")
			err = pg.generateFile(filename, path+string(filepath.Separator)+tmplname, Data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
