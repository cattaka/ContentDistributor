// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"

	// [START gae_go_env_data_imports]
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/api/option"
	// [END gae_go_env_data_imports]
	"context"
	"io"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/appengine/delay"
)

var (
	firebaseConfig = &firebase.Config{
					DatabaseURL:   "https://my-learning-project-140015.firebaseio.com",
					ProjectID:     "my-learning-project-140015",
					StorageBucket: "my-learning-project-140015.appspot.com",
	}
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [START gae_go_env_post_struct]
type Post struct {
        Author   string
        UserID   string
        Message  string
        Posted   time.Time
        ImageURL string
        Labels   []Label
}
// A Label is a description for a post's image.
type Label struct {
        Description string
        Score       float32
}

// [END gae_go_env_post_struct]

type templateParams struct {
	Notice string

	Name string
	// [START gae_go_env_template_params_fields]
	Message string

	Posts []Post
	// [END gae_go_env_template_params_fields]

}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

// labelFunc will be called asynchronously as a Cloud Task. labelFunc can
// be executed by calling labelFunc.Call(ctx, postID). If an error is returned
// the function will be retried.
var labelFunc = delay.Func("label-image", func(ctx context.Context, id int64) error {
        // Get the post to label.
        k := datastore.NewKey(ctx, "Post", "", id, nil)
        post := Post{}
        if err := datastore.Get(ctx, k, &post); err != nil {
                log.Errorf(ctx, "getting Post to label: %d, %v", id, err)
                return err
        }
        if post.ImageURL == "" {
                // Nothing to label.
                return nil
        }

        // Create a new vision client.
				opt := option.WithCredentialsFile("serviceAccountKey.json")
        client, err := vision.NewImageAnnotatorClient(ctx, opt)
        if err != nil {
                log.Errorf(ctx, "NewImageAnnotatorClient: %v", err)
                return err
        }
        defer client.Close()

        // Get the image and label it.
        image := vision.NewImageFromURI(post.ImageURL)
        labels, err := client.DetectLabels(ctx, image, nil, 5)
        if err != nil {
                log.Errorf(ctx, "Failed to detect labels: %v", err)
                return err
        }

        for _, l := range labels {
                post.Labels = append(post.Labels, Label{
                        Description: l.GetDescription(),
                        Score:       l.GetScore(),
                })
        }

        // Update the database with the new labels.
        if _, err := datastore.Put(ctx, k, &post); err != nil {
                log.Errorf(ctx, "Failed to update image: %v", err)
                return err
        }
        return nil
})

// uploadFileFromForm uploads a file if it's present in the "image" form field.
func uploadFileFromForm(ctx context.Context, r *http.Request) (url string, err error) {
        // Read the file from the form.
        f, fh, err := r.FormFile("image")
        if err == http.ErrMissingFile {
                return "", nil
        }
        if err != nil {
                return "", err
        }

        // Ensure the file is an image. http.DetectContentType only uses 512 bytes.
        buf := make([]byte, 512)
        if _, err := f.Read(buf); err != nil {
                return "", err
        }
        if contentType := http.DetectContentType(buf); !strings.HasPrefix(contentType, "image") {
                return "", fmt.Errorf("not an image: %s", contentType)
        }
        // Reset f so subsequent calls to Read start from the beginning of the file.
        f.Seek(0, 0)

        // Create a storage client.
				opt := option.WithCredentialsFile("serviceAccountKey.json")
        client, err := storage.NewClient(ctx, opt)
        if err != nil {
                return "", err
        }
        storageBucket := client.Bucket(firebaseConfig.StorageBucket)

        // Random filename, retaining existing extension.
        u, err := uuid.NewV4()
        if err != nil {
                return "", fmt.Errorf("generating UUID: %v", err)
        }
        name := u.String() + path.Ext(fh.Filename)

        w := storageBucket.Object(name).NewWriter(ctx)

        // Warning: storage.AllUsers gives public read access to anyone.
        w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
        w.ContentType = fh.Header.Get("Content-Type")

        // Entries are immutable, be aggressive about caching (1 day).
        w.CacheControl = "public, max-age=86400"

        if _, err := io.Copy(w, f); err != nil {
                w.CloseWithError(err)
                return "", err
        }
        if err := w.Close(); err != nil {
                return "", err
        }

        const publicURL = "https://storage.googleapis.com/%s/%s"
        return fmt.Sprintf(publicURL, firebaseConfig.StorageBucket, name), nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// [START gae_go_env_new_context]
	ctx := appengine.NewContext(r)
	// [END gae_go_env_new_context]
	params := templateParams{}

	// [START gae_go_env_new_query]
	q := datastore.NewQuery("Post").Order("-Posted").Limit(20)
	// [END gae_go_env_new_query]
	// [START gae_go_env_get_posts]
	if _, err := q.GetAll(ctx, &params.Posts); err != nil {
		log.Errorf(ctx, "Getting posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest posts. Refresh?"
		indexTemplate.Execute(w, params)
		return
	}
	// [END gae_go_env_get_posts]

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}

	message := r.FormValue("message")

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	// Create a new Firebase App.
//	app, err := firebase.NewApp(ctx, firebaseConfig)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
	        params.Notice = "Couldn't authenticate. Try logging in again?1"
	        params.Message = message // Preserve their message so they can try again.
	        indexTemplate.Execute(w, params)
	        return
	}
	// Create a new authenticator for the app.
	auth, err := app.Auth(ctx)
	if err != nil {
	        params.Notice = "Couldn't authenticate. Try logging in again?2"
	        params.Message = message // Preserve their message so they can try again.
	        indexTemplate.Execute(w, params)
	        return
	}
	// Verify the token passed in by the user is valid.
	tok, err := auth.VerifyIDTokenAndCheckRevoked(ctx, r.FormValue("token"))
	if err != nil {
	        params.Notice = "Couldn't authenticate. Try logging in again?3"
	        params.Message = message // Preserve their message so they can try again.
	        indexTemplate.Execute(w, params)
	        return
	}
	// Use the validated token to get the user's information.
	user, err := auth.GetUser(ctx, tok.UID)
	if err != nil {
	        params.Notice = "Couldn't authenticate. Try logging in again?4"
	        params.Message = message // Preserve their message so they can try again.
	        indexTemplate.Execute(w, params)
	        return
	}

	// It's a POST request, so handle the form submission.
	// [START gae_go_env_new_post]
	post := Post{
    UserID:  user.UID, // Include UserID in case Author isn't unique.
    Author:  user.DisplayName,
    Message: message,
    Posted:  time.Now(),
  }

	// [END gae_go_env_new_post]
	params.Name = post.Author

	if post.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}

	// Get the image if there is one.
	imageURL, err := uploadFileFromForm(ctx, r)
	if err != nil {
	        w.WriteHeader(http.StatusBadRequest)
	        params.Notice = "Error saving image: " + err.Error()
	        params.Message = post.Message // Preserve their message so they can try again.
	        indexTemplate.Execute(w, params)
	        return
	}
	post.ImageURL = imageURL

	key := datastore.NewIncompleteKey(ctx, "Post", nil)

	// [START gae_go_env_new_key]
	// [END gae_go_env_new_key]
	// [START gae_go_env_add_post]
	if key, err = datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// [END gae_go_env_add_post]

	// Only look for labels if the post has an image.
	if imageURL != "" {
	        // Run labelFunc. This will start a new Task in the background.
	        if err := labelFunc.Call(ctx, key.IntID()); err != nil {
	                log.Errorf(ctx, "delay Call %d, %v", key.IntID(), err)
	        }
	}

	// Prepend the post that was just added.
	// [START gae_go_env_prepend_post]
	params.Posts = append([]Post{post}, params.Posts...)
	// [END gae_go_env_prepend_post]

	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)

	indexTemplate.Execute(w, params)
}
