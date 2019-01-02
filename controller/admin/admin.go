package admin

import (
	"net/http"
	"html/template"
)

const PathPrefix = "/admin/"

type templateParams struct {
	Notice string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := template.Must(template.ParseFiles("template/admin/index.html"))
	params := templateParams{}
	htmlTemplate.Execute(w, params)
}
