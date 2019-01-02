package entity

import (
	"google.golang.org/appengine/datastore"
)

type Distribution struct {
	Key   *datastore.Key `datastore:"__key__"`
	Title string
	//ExpiredAt     date.Date
	//RealExpiredAt date.Date
	CoverImageURL string
}

type DistributionFile struct {
	Key      *datastore.Key `datastore:"__key__"`
	Parent   *datastore.Key
	FileName string
}

type DistributionCode struct {
	Key     *datastore.Key `datastore:"__key__"`
	IndexId string
	Parent  *datastore.Key
	Code    string
	Count   int
}
