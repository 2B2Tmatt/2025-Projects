package templates

import (
	"html/template"
	"log"
	"net/http"
)

func Render(w http.ResponseWriter, data any, filepath ...string) error {
	tmpl, err := template.ParseFiles(filepath...)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		log.Println(err)
		return err
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		log.Println(err)
		return err
	}
	return nil
}

// func Render(w http.ResponseWriter, data any, files ...string) error {
// 	tmpl, err := template.ParseFiles(files...)
// 	if err != nil {
// 		return err
// 	}
// 	return tmpl.Execute(w, data)
// }
