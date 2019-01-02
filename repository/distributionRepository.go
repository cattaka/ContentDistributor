package repository

import (
	"context"
	"github.com/cattaka/ContentDistributor/entity"
	"google.golang.org/appengine/datastore"
)

func FindDistributionsAll(ctx context.Context) ([]entity.Distribution, error) {
	var items []entity.Distribution
	q := datastore.NewQuery("Distribution").Order("-ExpiredAt").Order("Title")
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i< len(keys) && i< len(items);i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistribution(ctx context.Context, key *datastore.Key) (*entity.Distribution, error) {
	item := entity.Distribution{}
	err := datastore.Get(ctx, key, &item)
	if err == nil { item.Key = key }
	return &item, err
}

func FindDistributionFiles(ctx context.Context, parentKey *datastore.Key) ([]entity.DistributionFile, error) {
	var items []entity.DistributionFile
	q := datastore.NewQuery("DistributionFile").Filter("parentKey =", parentKey).Order("FileName")
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i< len(keys) && i< len(items);i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistributionCodes(ctx context.Context, parentKey *datastore.Key) ([]entity.DistributionCode, error) {
	var items []entity.DistributionCode
	q := datastore.NewQuery("DistributionCode").Filter("parentKey =", parentKey).Order("IndexId")
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i< len(keys) && i< len(items);i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func SaveDistribution(ctx context.Context, item *entity.Distribution) (*entity.Distribution, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "Distribution", nil)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func SaveDistributionFile(ctx context.Context, item *entity.DistributionFile) (*entity.DistributionFile, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionFile", nil)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func SaveDistributionCode(ctx context.Context, item *entity.DistributionCode) (*entity.DistributionCode, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionCode", nil)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}
