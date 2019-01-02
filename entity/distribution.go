package entity

import (
	"google.golang.org/appengine/datastore"
		"time"
)

type Distribution struct {
	Key   *datastore.Key `datastore:""`
	Title string
	ExpiredAt     time.Time
	RealExpiredAt time.Time
	CoverImageURL string
}

type DistributionFile struct {
	Key      *datastore.Key `datastore:""`
	Parent   *datastore.Key
	FileName string
}

type DistributionCode struct {
	Key     *datastore.Key `datastore:""`
	IndexId string
	Parent  *datastore.Key
	Code    string
	Count   int
}
