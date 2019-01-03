package storageUtil

import (
	"context"
	"github.com/cattaka/ContentDistributor/core"
	"mime/multipart"
	"io"
	"fmt"
	"cloud.google.com/go/storage"
)

func UploadFile(
	ctx context.Context,
	cb core.CoreBundle,
	f multipart.File,
	fh *multipart.FileHeader,
	fileFullPath string,
	makePublic bool) (url string, err error) {

	// Create a storage client.
	client, err := storage.NewClient(ctx, *(cb.ClientOption))
	if err != nil {
		return "", err
	}
	storageBucket := client.Bucket(cb.FirebaseConfig.StorageBucket)

	// Random filename, retaining existing extension.

	w := storageBucket.Object(fileFullPath).NewWriter(ctx)
	w.ContentType = fh.Header.Get("Content-Type")
	if makePublic {
		w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	}

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=300"

	if _, err := io.Copy(w, f); err != nil {
		w.CloseWithError(err)
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, cb.FirebaseConfig.StorageBucket, fileFullPath), nil
}
