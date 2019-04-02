// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"google.golang.org/appengine"
		"google.golang.org/api/option"
	"github.com/cattaka/ContentDistributor/core"
	"github.com/dsoprea/go-appengine-sessioncascade"
	"github.com/cattaka/ContentDistributor/router"
	"io/ioutil"
	"encoding/json"
	)

var (
	sessionSecret = []byte("SessionSecret")
	coreBundle = core.CoreBundle{}
)

func main() {
	var firebaseConfig core.FirebaseConfig
	if bytes, err := ioutil.ReadFile("firebaseConfig.json"); err != nil {
		panic(err)
	} else if err := json.Unmarshal(bytes, &firebaseConfig); err != nil {
		panic(err)
	}

	clientOption := option.WithCredentialsFile("serviceAccountKey.json")

	var aclConfig core.AclConfig
	if bytes, err := ioutil.ReadFile("aclConfig.json"); err != nil {
		panic(err)
	} else if err := json.Unmarshal(bytes, &aclConfig); err != nil {
		panic(err)
	}
	coreBundle = core.CoreBundle{
		SessionStore: cascadestore.NewCascadeStore(cascadestore.DistributedBackends, sessionSecret),
		ClientOption: &clientOption,
		FirebaseConfig: &firebaseConfig,
		AclConfig: &aclConfig,
	}

	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	router.IndexHandler(coreBundle, w, r)
}
