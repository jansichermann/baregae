package utils

import (
	"appengine"
	"appengine/datastore"
)

type DatastoreObject interface {
	ObjectId() string
	EntityType() string
}

func keyInContextForObject(ctx appengine.Context, obj DatastoreObject) datastore.Key {
	return datastore.NewKey(ctx, obj.EntityType(), obj.ObjectId(), 0, nil)
}

func PutObject(ctx appengine.Context, obj DatastoreObject) error {
	if _, err := datastore.Put(ctx, keyInContextForObject(ctx, obj), obj); err != nil {
		return err
	}
	return nil
}

func GetObject(ctx appengine.Context, obj DatastoreObject) error {
	key := datastore.NewKey(ctx, obj.EntityType(), obj.ObjectId(), 0, nil)
	if err := datastore.Get(ctx, keyInContextForObject(ctx, obj), obj); err != nil {
		return err
	}
	return nil
}

func DeleteObject(ctx appengine.Context, obj DatastoreObject) error {
	key := datastore.NewKey(ctx, obj.EntityType(), obj.ObjectId(), 0, nil)
	return datastore.Delete(ctx, keyInContextForObject(ctx, obj))
}
