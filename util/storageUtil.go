package main

import (
	"net/http"
	"strings"
	"fmt"
	"google.golang.org/api/option"
	"github.com/satori/go.uuid"
	"path"
	"io"
	"context"
	"google.golang.org/cloud/storage"
	"github.com/cattaka/ContentDistributor/core"
)

func uploadFileFromForm(ctx context.Context, cb core.CoreBundle, r *http.Request) (url string, err error) {
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
	storageBucket := client.Bucket(cb.ClientOption.StorageBucket)

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
