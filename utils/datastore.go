package utils

import (
	"appengine"
	"appengine/datastore"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

type DatastoreObject interface {
	ObjectId() string
	EntityType() string
}

type SettableId interface {
	SetObjectId(i string)
}

func keyInContextForObject(ctx appengine.Context, obj DatastoreObject) *datastore.Key {
	k := datastore.NewKey(ctx, obj.EntityType(), obj.ObjectId(), 0, nil)
	return k
}

func PutObject(ctx appengine.Context, obj DatastoreObject) error {
	if _, err := datastore.Put(ctx, keyInContextForObject(ctx, obj), obj); err != nil {
		return err
	}
	return nil
}

func GetObject(ctx appengine.Context, obj DatastoreObject) error {
	if err := datastore.Get(ctx, keyInContextForObject(ctx, obj), obj); err != nil {
		return err
	}
	return nil
}

func DeleteObject(ctx appengine.Context, obj DatastoreObject) error {
	return datastore.Delete(ctx, keyInContextForObject(ctx, obj))
}

/*
NewID ...
*/
func NewID() string {
	tb := make([]byte, 64)
	tc := binary.PutUvarint(tb, uint64(time.Now().UnixNano()))
	rb := make([]byte, 64)
	rc := binary.PutUvarint(rb, uint64(rand.Int63()))
	b := append(tb[:tc], rb[:rc]...)
	return fmt.Sprintf("%x", md5.Sum(b))
}
