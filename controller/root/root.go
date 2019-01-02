package root

import (
	"net/http"
	"html/template"
)

const PathPrefix = "/"

type templateParams struct {
	Notice string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := template.Must(template.ParseFiles("template/index.html"))
	params := templateParams{}
	htmlTemplate.Execute(w, params)
}
