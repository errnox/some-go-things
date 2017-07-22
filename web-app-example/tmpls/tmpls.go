package tmpls

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/<username>/misc/routes"
)

const (
	templateFileExtension  = ".html.tmpl"
	templatesDirName       = "templates"
	commonTemplatesDirName = "common"
	commonTemplatesPrefix  = "common"
)

type ParsedHTML struct {
	HTML string
}

func (h *ParsedHTML) Write(p []byte) (n int, err error) {
	for _, i := range p {
		h.HTML = fmt.Sprintf("%s%s", h.HTML, string(i))
	}
	return len(p), nil
}

var htmlTemplates *template.Template
var currentPath string
var Flash string
var router *mux.Router

func init() {
	// htmlTemplates = parseHTMLTemplates()
}

func SetRouter(r *mux.Router) {
	router = r
}

func Get(s string) *template.Template {
	return GetAll().Lookup(s + templateFileExtension)
}

func addFuncs(t *template.Template) {
	t.Funcs(template.FuncMap{
		"linkTo":      linkToFunc,
		"activeClass": activeClassFunc,
		"flash":       flashFunc,
		"url":         urlFunc,
		"add":         addFunc,
		"form":        formFunc,
	})
}

func Render(w io.Writer, data interface{}, name, layout string) {
	currentPath = name

	templatePath := fmt.Sprintf("./%s/%s%s", templatesDirName, name,
		templateFileExtension)
	layoutPath := fmt.Sprintf("./%s/%s%s", templatesDirName, layout,
		templateFileExtension)

	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Could not read file %s: %s\n", templatePath, err)
	}
	layoutContent, err := ioutil.ReadFile(layoutPath)
	if err != nil {
		fmt.Printf("Could not read file %s: %s\n", layoutPath, err)
	}
	commonContent := readCommonTemplates()
	content := fmt.Sprintf("%s\n%s\n\n%s", commonContent,
		templateContent, layoutContent)

	temp := template.New("")
	addFuncs(temp)
	temp.Parse(content)
	if err != nil {
		fmt.Println("Could not parse template: ", err)
	}
	temp.Execute(w, data)
	Flash = ""
}

func GetAll() *template.Template {
	return htmlTemplates
}

func parseHTMLTemplates() *template.Template {
	t, err := template.ParseGlob("./templates/**.html.tmpl")
	if err != nil {
		fmt.Println("Could not parse template: ", err)
	}
	return t
}

func readCommonTemplates() string {
	s := ""
	glob := fmt.Sprintf("./%s/%s/**%s",
		templatesDirName, commonTemplatesDirName,
		templateFileExtension)
	files, err := filepath.Glob(glob)
	if err != nil {
		fmt.Printf("Could not list files in directory matching glob %s: %s",
			glob, err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Could not read file %s: %s", file,
				err)
		}
		s = fmt.Sprintf("%s\n%s", s, content)
	}
	return s
}

// FuncMap functions

func linkToFunc(s string) template.HTML {
	return template.HTML(routes.Routes[s])
}

func urlFunc(name string, pairs ...interface{}) template.HTML {
	args := []string{}
	for _, x := range pairs {
		args = append(args, fmt.Sprintf("%v", x))
	}
	route, err := router.Get(name).URL(args...)
	if err != nil {
		log.Fatal(err)
	}
	url := route.Path
	return template.HTML(url)
}

func addFunc(args ...int) template.HTML {
	x := 0
	for _, n := range args {
		x = x + n
	}
	return template.HTML(fmt.Sprintf("%d", x))
}

func formFunc(x interface{}) template.HTML {
	s := "<form>"
	o := reflect.ValueOf(x).Elem()
	for i := 0; i < o.NumField(); i++ {
		if o.Field(i).Kind() == reflect.Int {
			s += fmt.Sprintf(`<div class="form-group">
  <label>%v</label>
  <input type="number" class="form-control" value="%v"/>
</div>`, o.Type().Field(i).Tag.Get("name"), o.Field(i).Interface())
		} else {
			s += fmt.Sprintf(`<div class="form-group">
  <label>%v</label>
  <input type="text" class="form-control" value="%v"/>
</div>`, o.Type().Field(i).Tag.Get("name"), o.Field(i).Interface())
		}
	}

	s += "</form>"
	return template.HTML(s)
}

func activeClassFunc(path string, classes ...string) template.HTML {
	s := ""
	if path == currentPath {
		for _, c := range classes {
			s = fmt.Sprintf("%s %s", s, strings.TrimSpace(c))
		}
	}
	s = strings.TrimSpace(s)
	return template.HTML(s)
}

func flashFunc() template.HTML {
	f := fmt.Sprintf("%s", Flash)
	return template.HTML(f)
}
