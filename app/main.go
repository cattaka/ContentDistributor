// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
		"net/http"

	"google.golang.org/appengine"
	"github.com/cattaka/ContentDistributor/router"
)

func main() {
	http.HandleFunc("/", router.IndexHandler)
	appengine.Main()
}
