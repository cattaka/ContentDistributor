package router

import (
	"net/http"
	"strings"
	"github.com/cattaka/ContentDistributor/controller/admin"
	"github.com/cattaka/ContentDistributor/controller/root"
	"github.com/cattaka/ContentDistributor/core"
	"github.com/dsoprea/go-appengine-sessioncascade"
	"google.golang.org/api/option"
)

var (
	sessionSecret = []byte("SessionSecret")
	clientOption = option.WithCredentialsFile("serviceAccountKey.json")
	coreBundle = core.CoreBundle{
		SessionStore: cascadestore.NewCascadeStore(cascadestore.DistributedBackends, sessionSecret),
		ClientOption: &clientOption,
	}
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, admin.PathPrefix) {
		admin.IndexHandler(coreBundle,  w, r)
		return
	} else if strings.HasPrefix(r.URL.Path, root.PathPrefix) {
		root.IndexHandler(coreBundle, w, r)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
