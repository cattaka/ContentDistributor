package core

import (
	"github.com/dsoprea/go-appengine-sessioncascade"
	"github.com/gorilla/sessions"
	"google.golang.org/api/option"
)

var (
	sessionSecret = []byte("SessionSecret")
	sessionStore = cascadestore.NewCascadeStore(cascadestore.DistributedBackends, sessionSecret)
)

type CoreBundle struct {
	SessionStore *cascadestore.CascadeStore
	Session *sessions.Session
	ClientOption *option.ClientOption
}