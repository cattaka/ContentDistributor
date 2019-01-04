package repository

import (
	"context"
	"github.com/cattaka/ContentDistributor/entity"
	"google.golang.org/appengine/datastore"
	"github.com/cattaka/ContentDistributor/util"
)

func FindDistributionsAll(ctx context.Context, withDisabled bool) ([]entity.Distribution, error) {
	var items []entity.Distribution
	q := datastore.NewQuery("Distribution").Order("-ExpiredAt").Order("Title")
	if !withDisabled {
		q = q.Filter("Disabled =", false)
	}
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i < len(keys) && i < len(items); i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistribution(ctx context.Context, key *datastore.Key) (*entity.Distribution, error) {
	item := entity.Distribution{}
	err := datastore.Get(ctx, key, &item)
	if err == nil {
		item.Key = key
	}
	return &item, err
}

func FindDistributionFile(ctx context.Context, key *datastore.Key) (*entity.DistributionFile, error) {
	item := entity.DistributionFile{}
	err := datastore.Get(ctx, key, &item)
	if err == nil {
		item.Key = key
	}
	return &item, err
}

func FindDistributionFiles(ctx context.Context, parentKey *datastore.Key, withDisabled bool) ([]entity.DistributionFile, error) {
	var items []entity.DistributionFile
	q := datastore.NewQuery("DistributionFile").Filter("Parent =", parentKey).Order("FileName")
	if !withDisabled {
		q = q.Filter("Disabled =", false)
	}
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i < len(keys) && i < len(items); i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistributionCode(ctx context.Context, key *datastore.Key) (*entity.DistributionCode, error) {
	item := entity.DistributionCode{}
	err := datastore.Get(ctx, key, &item)
	if err == nil {
		item.Key = key
	}
	return &item, err
}

func FindDistributionCodeByCode(ctx context.Context, code string) ([]entity.DistributionCode, error) {
	var items []entity.DistributionCode
	q := datastore.NewQuery("DistributionCode").Filter("Code =", code).Order("IdLabel")
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i < len(keys) && i < len(items); i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistributionCodes(ctx context.Context, parentKey *datastore.Key, withDisabled bool) ([]entity.DistributionCode, error) {
	var items []entity.DistributionCode
	q := datastore.NewQuery("DistributionCode").Filter("Parent =", parentKey).Order("IdLabel")
	if !withDisabled {
		q = q.Filter("Disabled =", false)
	}
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i < len(keys) && i < len(items); i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistributionCodesByTag(ctx context.Context, tagKey *datastore.Key, withDisabled bool) ([]entity.DistributionCode, error) {
	var items []entity.DistributionCode
	q := datastore.NewQuery("DistributionCode").Filter("GenerationTag =", tagKey).Order("IdLabel")
	if !withDisabled {
		q = q.Filter("Disabled =", false)
	}
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i < len(keys) && i < len(items); i++ {
			items[i].Key = keys[i]
		}
	}
	return items, err
}

func FindDistributionGenerationTag(ctx context.Context, key *datastore.Key) (*entity.DistributionGenerationTag, error) {
	item := entity.DistributionGenerationTag{}
	err := datastore.Get(ctx, key, &item)
	if err == nil {
		item.Key = key
	}
	return &item, err
}

func FindDistributionGenerationTags(ctx context.Context, parentKey *datastore.Key, withDisabled bool) ([]entity.DistributionGenerationTag, error) {
	var items []entity.DistributionGenerationTag
	q := datastore.NewQuery("DistributionGenerationTag").Filter("Parent =", parentKey).Order("Name")
	if !withDisabled {
		q = q.Filter("Disabled =", false)
	}
	keys, err := q.GetAll(ctx, &items)
	if err == nil {
		for i := 0; i < len(keys) && i < len(items); i++ {
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
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionFile", item.Parent)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func SaveDistributionCode(ctx context.Context, item *entity.DistributionCode) (*entity.DistributionCode, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionCode", item.Parent)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func SaveDistributionCodes(ctx context.Context, items *[]entity.DistributionCode) error {
	var tempKeys = make([]*datastore.Key, len(*items))
	for i, v := range *items {
		if v.Key == nil {
			tempKeys[i] = datastore.NewIncompleteKey(ctx, "DistributionCode", v.Parent)
		} else {
			tempKeys[i] = v.Key
		}
	}
	if keys, err := datastore.PutMulti(ctx, tempKeys, *items); err == nil {
		for i := 0; i < len(keys) && i < len(*items); i++ {
			(*items)[i].Key = keys[i]
		}
		return err
	} else {
		return err
	}

}

func SaveDistributionGenerationTag(ctx context.Context, item *entity.DistributionGenerationTag) (*entity.DistributionGenerationTag, error) {
	var err error
	if item.Key == nil {
		item.Key = datastore.NewIncompleteKey(ctx, "DistributionGenerationTag", item.Parent)
	}
	item.Key, err = datastore.Put(ctx, item.Key, item)
	return item, err
}

func GenNextUniqueCode(ctx context.Context, length int) (string, error) {
	var code string
	err := datastore.RunInTransaction(ctx, func(c context.Context) error {
		for ; ; {
			code = util.RandString(length)
			nextKey := datastore.NewKey(ctx, "UniqueCode", code, 0, nil)
			var item = entity.UniqueCode{}
			if err := datastore.Get(c, nextKey, &item); err == nil {
				continue // Already exists
			} else if err != datastore.ErrNoSuchEntity {
				return err
			}
			item.Code = code
			if _, err := datastore.Put(c, nextKey, &item); err != nil {
				return err
			}
			break
		}
		return nil
	}, nil)
	return code, err
}
