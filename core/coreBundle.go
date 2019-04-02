package core

import (
	"github.com/dsoprea/go-appengine-sessioncascade"
	"github.com/gorilla/sessions"
	"google.golang.org/api/option"
)

var (
	sessionSecret = []byte("SessionSecret")
	sessionStore  = cascadestore.NewCascadeStore(cascadestore.DistributedBackends, sessionSecret)
)

type FirebaseConfig struct {
	ApiKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	DatabaseURL       string `json:"databaseURL"`
	ProjectId         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderId string `json:"messagingSenderId"`
}

type AclConfig struct {
	AvailableAccounts []string `json:"availableAccounts"`
}

type CoreBundle struct {
	SessionStore   *cascadestore.CascadeStore
	Session        *sessions.Session
	ClientOption   *option.ClientOption
	FirebaseConfig *FirebaseConfig
	AclConfig      *AclConfig
}
