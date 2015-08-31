package utils

import (
	"appengine"
	"appengine/datastore"
)

type DatastoreObject interface {
	ObjectId() string
	EntityType() string
}

func PutObject(ctx appengine.Context, obj DatastoreObject) error {
	key := datastore.NewKey(ctx, obj.EntityType(), obj.ObjectId(), 0, nil)
	if _, err := datastore.Put(ctx, key, obj); err != nil {
		return err
	}
	return nil
}

func GetObject(ctx appengine.Context, obj DatastoreObject) error {
	key := datastore.NewKey(ctx, obj.EntityType(), obj.ObjectId(), 0, nil)
	if err := datastore.Get(ctx, key, obj); err != nil {
		return err
	}
	return nil
}
