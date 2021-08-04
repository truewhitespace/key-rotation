package rotation

import (
	"context"
	"time"
)

//InvalidTime produces a common reference time to indicate a key does not have a valid creation time or otherwise is
//not suitable for operations
func InvalidTime() time.Time {
	return time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)
}

//Key is a bridge to the underlying implementation of a particular KeyStore.
type Key interface {
	//Created provides the time the key was created.  If the key is otherwise not valid or inactive the key should provide
	//rotation.InvalidTime() as the result of this invocation.
	Created() time.Time
}

//KeyList is an anemic type reference for a set of keys...probably should add behavior or get rid of it.
type KeyList []Key

//KeyStore abstracts operations to be performed against a key store for rotational capabilities.  These are typically
//bound to a specific agent context such as a user or application.
type KeyStore interface {

	//CreateKey creates and persists a new key within the key store.  Key should be castable to the implementing type of
	//the key store for extraction of the specific credentials.
	CreateKey(ctx context.Context) (Key, error)

	//DeleteKey causes the destruction of a given key from the store.  The key should no longer be operable against the
	//target systems after the completion of this method.
	DeleteKey(ctx context.Context, key Key) error

	//ListKeys queries the target key store for
	ListKeys(ctx context.Context) (KeyList, error)

	//MaximumKeys is the total number of keys which can be created per user within this store.  This is useful in cases
	//where multiple keys maybe in grace, among others.
	MaximumKeys() int
}
