package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/doantaydo/Learning-GO_Web-Application/Basic-Web-Application/packages/config"
	"github.com/doantaydo/Learning-GO_Web-Application/Basic-Web-Application/packages/models"
)

var template_cache = make(map[string]*template.Template)

var app *config.AppConfig

// NewTemplates sets the config for the package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// get the template cache from AppConfig
	var templateSet map[string]*template.Template
	if app.UseCache {
		templateSet = app.TemplateCache
	} else {
		templateSet, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := templateSet[tmpl]
	if ok == false {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	// render template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all file *.page.tmpl in templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	// range through all file ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		parTmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			parTmpl, err = parTmpl.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = parTmpl
	}

	return myCache, nil
}

func RenderTemplate2(w http.ResponseWriter, tmpl string) {
	fmt.Println("Check " + tmpl)

	var parTmpl *template.Template
	var err error

	// check if has template in cache
	_, inMap := template_cache[tmpl]

	if !inMap {
		// need to create template in cache
		err = CreateTemplateCache2(tmpl)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Create template " + tmpl + " in cache.")
	} else {
		// have template in cache
		fmt.Println("Using template " + tmpl + " in cache.")
	}

	parTmpl = template_cache[tmpl]
	err = parTmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateTemplateCache2(tmpl string) error {
	newTemplate, err := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")

	if err != nil {
		return err
	}

	template_cache[tmpl] = newTemplate
	return nil
}
