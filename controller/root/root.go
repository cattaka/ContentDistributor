package root

import (
	"net/http"
		"github.com/cattaka/ContentDistributor/core"
	"context"
	"strings"
	"google.golang.org/appengine"
	"html/template"
		"github.com/cattaka/ContentDistributor/repository"
	"github.com/cattaka/ContentDistributor/entity"
	"fmt"
	"cloud.google.com/go/storage"
	"github.com/cattaka/ContentDistributor/util"
)

const PathPrefix = "/"

type templateParams struct {
	Notice string
}

func IndexHandler(cb core.CoreBundle, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if download(&ctx, cb, w, r) {
		// ok
	} else {
		htmlTemplate := template.Must(template.ParseFiles("template/index.html"))
		params := templateParams{}
		htmlTemplate.Execute(w, params)
	}
}


func download(ctx *context.Context, cb core.CoreBundle, w http.ResponseWriter, r *http.Request) bool {
	segments := strings.Split(r.URL.Path, "/")

	if len(segments) < 3 {
		return false
	}
	var keyStr = segments[1]
	var fileName = segments[2]

	var distribution *entity.Distribution
	var code *entity.DistributionCode
	var files []entity.DistributionFile
	if items, err := repository.FindDistributionCodeByCode(*ctx, keyStr); err != nil || len(items) == 0 {
		return false
	} else if d, err := repository.FindDistribution(*ctx, items[0].Parent); err != nil {
		return false
	} else if f, err := repository.FindDistributionFiles(*ctx, items[0].Parent, false); err != nil {
		return false
	} else {
		code = &items[0]
		distribution = d
		files = f
	}

	var file *entity.DistributionFile = nil
	for _,f := range files {
		if f.FileName == fileName {
			file = &f
			break
		}
	}
	if file == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not found"))
		return true
	}

	downloadFileFullPath := fmt.Sprintf("dist/%s/%s", code.Code, fileName)
	origFileFullPath := fmt.Sprintf("orig/%s/%s", distribution.Key.Encode(), fileName)

	client, err := storage.NewClient(*ctx, *(cb.ClientOption))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return true
	}
	storageBucket := client.Bucket(cb.FirebaseConfig.StorageBucket)
	origObject := storageBucket.Object(origFileFullPath)
	downloadObj := storageBucket.Object(downloadFileFullPath)
	var acls []storage.ACLRule
	if a, err := downloadObj.ACL().List(*ctx); err == nil {
		acls = a
	} else if _, err := downloadObj.CopierFrom(origObject).Run(*ctx); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return true
	} else if a, err := downloadObj.ACL().List(*ctx); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return true
	} else {
		acls = a
	}

	var readable = false
	for _,acl := range acls {
		if acl.Entity == storage.AllUsers && acl.Role == storage.RoleReader {
			readable = true
			break
		}
	}
	if !readable {
		if err := downloadObj.ACL().Set(*ctx, storage.AllUsers, storage.RoleReader);err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return true
		}
	}

	code.Count++
	repository.SaveDistributionCode(*ctx, code)

	http.Redirect(w, r, fmt.Sprintf(util.STORAGE_URL_FORMAT, cb.FirebaseConfig.StorageBucket, downloadFileFullPath), http.StatusFound)
	return true
}