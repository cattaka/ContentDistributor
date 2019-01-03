package admin

import (
	"net/http"
	"html/template"
	"google.golang.org/appengine"
	"github.com/cattaka/ContentDistributor/core"
	"context"
	"github.com/cattaka/ContentDistributor/entity"
	"github.com/cattaka/ContentDistributor/repository"
	"google.golang.org/appengine/datastore"
	"time"
	"fmt"
	"firebase.google.com/go"
	"github.com/cattaka/ContentDistributor/util"
	"regexp"
	"strconv"
)

const (
	sessionName  = "MainSession"
	PathPrefix   = "/admin/"
	KeyAuthToken = "AuthToken"
)

type templateParams struct {
	Notice            string
	SignedIn          bool
	Distributions     []entity.Distribution
	Distribution      *entity.Distribution
	DistributionFiles []entity.DistributionFile
	FirebaseConfig    core.FirebaseConfig
	DistributionCodes []entity.DistributionCode
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
	} else if r.Method == "GET" && r.URL.Path == PathPrefix+"signInOut" {
		showSignInOut(&ctx, cb, w, r)
	} else if r.Method == "GET" && r.URL.Path == PathPrefix+"editDistribution" {
		showEditDistribution(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"editDistribution" {
		postEditDistribution(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"addDistributionFile" {
		addDistributionFile(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"deleteDistributionFile" {
		deleteDistributionFile(&ctx, cb, w, r)
	} else if r.Method == "GET" && r.URL.Path == PathPrefix+"editDistributionCodes" {
		showEditDistributionCodes(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"generateDistributionCodes" {
		generateDistributionCodes(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"signIn" {
		signIn(&ctx, cb, w, r)
	} else if r.Method == "POST" && r.URL.Path == PathPrefix+"signOut" {
		signOut(&ctx, cb, w, r)
	}
}

func showIndex(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	params := templateParams{}
	if _, found := cb.Session.Values[KeyAuthToken]; found {
		params.SignedIn = true
		params.Distributions, _ = repository.FindDistributionsAll(*ctx, false)
	}

	htmlTemplate := template.Must(template.ParseFiles("template/admin/index.html"))
	htmlTemplate.Execute(w, params)
}

func showSignInOut(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	params := templateParams{}
	if _, found := cb.Session.Values[KeyAuthToken]; found {
		params.SignedIn = true
	}
	params.FirebaseConfig = *cb.FirebaseConfig

	htmlTemplate := template.Must(template.ParseFiles("template/admin/signInOut.html"))
	htmlTemplate.Execute(w, params)
}

func signIn(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	app, err := firebase.NewApp(*ctx, nil, *cb.ClientOption)
	if err != nil {
		panic(err)
	}
	auth, err := app.Auth(*ctx)
	if err != nil {
		panic(err)
	}
	tok, err := auth.VerifyIDTokenAndCheckRevoked(*ctx, r.FormValue("token"))
	if err != nil {
		panic(err)
	}
	_, err = auth.GetUser(*ctx, tok.UID)
	if err != nil {
		panic(err)
	}

	cb.Session.Values[KeyAuthToken] = token
	cb.Session.Save(r, w)

	http.Redirect(w, r, PathPrefix, http.StatusFound)
}

func signOut(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	if _, found := cb.Session.Values[KeyAuthToken]; found {
		delete(cb.Session.Values, KeyAuthToken)
		cb.Session.Save(r, w)
	}
	http.Redirect(w, r, PathPrefix, http.StatusFound)
}

func showEditDistribution(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	params := templateParams{}
	if _, found := cb.Session.Values[KeyAuthToken]; !found {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	}
	params.SignedIn = true

	if k, err := datastore.DecodeKey(r.FormValue("Key")); err == nil {
		if item, e2 := repository.FindDistribution(*ctx, k); e2 == nil {
			params.Distribution = item
			if files, e3 := repository.FindDistributionFiles(*ctx, k, false); e3 == nil {
				params.DistributionFiles = files
			}
		}
	}
	if params.Distribution == nil {
		item := entity.Distribution{}
		params.Distribution = &item
	}

	htmlTemplate := template.Must(template.ParseFiles("template/admin/editDistribution.html"))
	htmlTemplate.Execute(w, params)
}

func postEditDistribution(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	if _, found := cb.Session.Values[KeyAuthToken]; !found {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	}

	format := "2006-01-02"
	k, _ := datastore.DecodeKey(r.FormValue("Key"))
	expiredAt, _ := time.Parse(format, r.FormValue("ExpiredAt"))
	realExpiredAt, _ := time.Parse(format, r.FormValue("RealExpiredAt"))
	item := entity.Distribution{
		Key:           k,
		Title:         r.FormValue("Title"),
		ExpiredAt:     expiredAt,
		RealExpiredAt: realExpiredAt,
		CoverImageUrl: r.FormValue("CoverImageUrl"),
	}
	repository.SaveDistribution(*ctx, &item)

	http.Redirect(w, r, fmt.Sprintf("%seditDistribution?Key=%s", PathPrefix, item.Key.Encode()), http.StatusFound)
}

func addDistributionFile(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	if _, found := cb.Session.Values[KeyAuthToken]; !found {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	}

	var key *datastore.Key
	if k, err := datastore.DecodeKey(r.FormValue("Key")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if _, err := repository.FindDistribution(*ctx, k); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {
		key = k
	}

	fileNameRegex := "^[a-zA-Z0-9][a-zA-Z0-9\\-_\\.]*$"
	fileName := r.FormValue("FileName")
	rg := regexp.MustCompile(fileNameRegex)
	if !rg.Match([]byte(fileName)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("file name must match `%s`", fileNameRegex)))
		return
	}

	fileFullPath := fmt.Sprintf("o/%s/%s", key.Encode(), fileName)

	var url string
	if f, fh, err := r.FormFile("File"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if u, err := storageUtil.UploadFile(*ctx, cb, f, fh, fileFullPath, false); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {
		url = u
	}

	distributionFile := entity.DistributionFile{Parent: key, FileName: fileName, Url: url}
	if _, err := repository.SaveDistributionFile(*ctx, &distributionFile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%seditDistribution?Key=%s", PathPrefix, key.Encode()), http.StatusFound)
}

func deleteDistributionFile(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	var df *entity.DistributionFile
	if k, err := datastore.DecodeKey(r.FormValue("Key")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if i, err := repository.FindDistributionFile(*ctx, k); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {
		df = i
	}

	df.Disabled = true
	if _, err := repository.SaveDistributionFile(*ctx, df); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%seditDistribution?Key=%s", PathPrefix, df.Parent.Encode()), http.StatusFound)
}

func showEditDistributionCodes(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	params := templateParams{}
	if _, found := cb.Session.Values[KeyAuthToken]; !found {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	}
	params.SignedIn = true

	if k, err := datastore.DecodeKey(r.FormValue("Key")); err != nil {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	} else if item, e2 := repository.FindDistribution(*ctx, k); e2 != nil {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	} else if codes, e3 := repository.FindDistributionCodes(*ctx, k, false); e3 == nil {
		params.Distribution = item
		params.DistributionCodes = codes
	}

	htmlTemplate := template.Must(template.ParseFiles("template/admin/editDistributionCodes.html"))
	htmlTemplate.Execute(w, params)
}

func generateDistributionCodes(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	if _, found := cb.Session.Values[KeyAuthToken]; !found {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	}

	var distribution *entity.Distribution
	if k, err := datastore.DecodeKey(r.FormValue("Key")); err != nil {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	} else if item, e2 := repository.FindDistribution(*ctx, k); e2 != nil {
		http.Redirect(w, r, PathPrefix, http.StatusFound)
		return
	} else {
		distribution = item
	}

	var idFrom, idTo int
	if i, err := strconv.Atoi(r.FormValue("IdFrom")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if j, err := strconv.Atoi(r.FormValue("IdTo")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {
		idFrom = i
		idTo = j
	}

	idFormat := r.FormValue("IdFormat")

	var distributionCodes []entity.DistributionCode
	for i := idFrom; i <= idTo; i++ {
		distributionCodes = append(distributionCodes,
			entity.DistributionCode{
				Parent:   distribution.Key,
				IdLabel:  fmt.Sprintf(idFormat, i),
				Count:    0,
				Disabled: false,
			})
	}
	if err := repository.SaveDistributionCodes(*ctx, &distributionCodes) ; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%seditDistributionCodes?Key=%s", PathPrefix, distribution.Key.Encode()), http.StatusFound)
}
