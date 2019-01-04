package entity

import (
	"google.golang.org/appengine/datastore"
	"time"
)

type Distribution struct {
	Key           *datastore.Key `datastore:"-"`
	Title         string
	ExpiredAt     time.Time
	RealExpiredAt time.Time
	Contact       string
	CoverImageUrl string
	Disabled      bool
}

type DistributionFile struct {
	Key        *datastore.Key `datastore:"-"`
	Parent     *datastore.Key
	FileName   string
	ShortLabel string
	Url        string
	Disabled   bool
}

type DistributionCode struct {
	Key           *datastore.Key `datastore:"-"`
	Parent        *datastore.Key
	Code          string
	GenerationTag *datastore.Key
	IdLabel       string
	Count         int
	Disabled      bool
}

type DistributionGenerationTag struct {
	Key      *datastore.Key `datastore:"-"`
	Parent   *datastore.Key
	Name     string
	IdFormat string
	IdFrom   int
	IdTo     int
	Disabled bool
}

type UniqueCode struct {
	Key  *datastore.Key `datastore:"-"`
	Code string
}
