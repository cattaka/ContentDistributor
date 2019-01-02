package admin

import (
	"net/http"
	"html/template"
	"google.golang.org/appengine"
	"github.com/cattaka/ContentDistributor/core"
	"firebase.google.com/go"
	"context"
)

const (
	sessionName = "MainSession"
	PathPrefix  = "/admin/"
	KeyAuthToken = "AuthToken"
)

type templateParams struct {
	Notice string
}

func IndexHandler(cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	session, err := cb.SessionStore.Get(r, sessionName)
	if err != nil {
		panic(err)
	}
	cb.Session = session

	if r.Method == "GET" && r.URL.Path == PathPrefix {
		showIndex(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"signIn" {
		signIn(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"signOut" {
		signOut(&ctx, cb, w, r)
	}
}

func showIndex(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	htmlTemplate := template.Must(template.ParseFiles("template/admin/index.html"))
	params := templateParams{}
	htmlTemplate.Execute(w, params)
}

func signIn(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	app, err := firebase.NewApp(*ctx, nil, *cb.ClientOption)
	if err != nil { panic(err) }
	auth, err := app.Auth(*ctx)
	if err != nil { panic(err) }
	tok, err := auth.VerifyIDTokenAndCheckRevoked(*ctx, r.FormValue("token"))
	if err != nil { panic(err) }
	_, err = auth.GetUser(*ctx, tok.UID)
	if err != nil { panic(err) }

	cb.Session.Values[KeyAuthToken] = token

	http.Redirect(w, r, "/admin/", http.StatusFound)
}

func signOut(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	if _, found := cb.Session.Values[KeyAuthToken]; found {
		delete(cb.Session.Values, KeyAuthToken)
	}
	http.Redirect(w, r, "/admin/", http.StatusFound)
}
