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
	bytes, err := ioutil.ReadFile("firebaseConfig.json")
	if err != nil { panic(err) }
	var firebaseConfig core.FirebaseConfig
	if err := json.Unmarshal(bytes, &firebaseConfig); err != nil { panic(err) }

	clientOption := option.WithCredentialsFile("serviceAccountKey.json")
	coreBundle = core.CoreBundle{
		SessionStore: cascadestore.NewCascadeStore(cascadestore.DistributedBackends, sessionSecret),
		ClientOption: &clientOption,
		FirebaseConfig: &firebaseConfig,
	}

	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	router.IndexHandler(coreBundle, w, r)
}
