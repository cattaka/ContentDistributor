package repository

import (
	"context"
	"github.com/cattaka/ContentDistributor/entity"
	"google.golang.org/appengine/datastore"
)

func findDistributionsAll(ctx context.Context) ([]entity.Distribution, error) {
	var items []entity.Distribution
	q := datastore.NewQuery("Distribution").Order("-ExpiredAt").Order("Title")
	_, err := q.GetAll(ctx, &items)
	return items, err
}

func findDistribution(ctx context.Context, key *datastore.Key) (*entity.Distribution, error) {
	item := entity.Distribution{}
	err := datastore.Get(ctx, key, &item)
	return &item, err
}

func DistributionFiles(ctx context.Context, parentKey *datastore.Key) ([]entity.DistributionFile, error) {
	var items []entity.DistributionFile
	q := datastore.NewQuery("DistributionFile").Filter("parentKey =", parentKey).Order("FileName")
	_, err := q.GetAll(ctx, &items)
	return items, err
}

func DistributionCodes(ctx context.Context, parentKey *datastore.Key) ([]entity.DistributionCode, error) {
	var items []entity.DistributionCode
	q := datastore.NewQuery("DistributionCode").Filter("parentKey =", parentKey).Order("IndexId")
	_, err := q.GetAll(ctx, &items)
	return items, err
}

func saveDistribution(ctx context.Context, item *entity.Distribution) (*entity.Distribution, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "Distribution", nil)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func saveDistributionFile(ctx context.Context, item *entity.DistributionFile) (*entity.Distribution, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionFile", nil)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func saveDistributionCode(ctx context.Context, item *entity.DistributionCode) (*entity.Distribution, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionCode", nil)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}
