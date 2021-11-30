package env

import (
	"io/ioutil"
	"log"
	"text/template"
)

func ParseTemplates() map[string]*template.Template {
	var result = make(map[string]*template.Template)
	global, pages := GetTemplatesPathsAndNames()

	for i := 0; i < len(pages); i++ {
		var templ *template.Template
		var page string = "./templates/" + pages[i]
		var templatePathsArr []string = append([]string{page}, global...)

		templ, err := template.ParseFiles(templatePathsArr...)
		if err != nil {
			panic(err)
		}

		result[pages[i]] = templ
	}

	return result
}

func GetTemplatesPathsAndNames() ([]string, []string) {
	var pages []string
	var global []string

	root := "./templates"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			pages = append(pages, f.Name())
		}
	}

	root = "./templates/global"
	files, err = ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		global = append(global, "./templates/global/"+f.Name())
	}

	return global, pages
}
