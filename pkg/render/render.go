package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/clauseisenhardt/bookings/pkg/config"
	"github.com/clauseisenhardt/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate renders templates using http/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// tc, err := CreateTemplateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// for _, t := range tc {
	// 	fmt.Println("template: ", t)
	// }
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser:", err)
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	// err = parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error parsing template:", err)
	// 	return
	// }
	// fmt.Println("Error writing template to browser - take2:", err)
}

// CreateTemplateCache creates a template cach as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	fmt.Println("CreateTemplateCache...")
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		fmt.Println("Error getting templates:", err)
		return myCache, err
	}

	//fmt.Println("For pages...")
	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("Page is currently:", page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return myCache, err
		}
		//fmt.Println("Match template pages...")
		matches, err2 := filepath.Glob("./templates/*.layout.tmpl")
		if err2 != nil {
			fmt.Println("Error matching layout templates:", err2)
			return myCache, err
		}
		//fmt.Println("Templates pages matched...")
		if len(matches) > 0 {
			//fmt.Println("Templates pages matched > 0 ...")
			ts2, err3 := ts.ParseGlob("./templates/*.layout.tmpl")
			if err3 != nil {
				fmt.Println("Error parsing layout templates:", err3)
				return myCache, err
			}
			myCache[name] = ts2
		} else {
			fmt.Println("Error: no matching templates found!")
		}

	}

	return myCache, nil
}
