package root

import (
	"net/http"
	"html/template"
	"github.com/cattaka/ContentDistributor/core"
)

const PathPrefix = "/"

type templateParams struct {
	Notice string
}

func IndexHandler(cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	htmlTemplate := template.Must(template.ParseFiles("template/index.html"))
	params := templateParams{}
	htmlTemplate.Execute(w, params)
}
