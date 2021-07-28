package rotation

import (
	"context"
	"time"
)

func InvalidTime() time.Time {
	return time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)
}

//Key is a bridge to the underlying hierarchy of a particular KeyStore.
type Key interface {
	//Created provides the time the key was created.  If the key is otherwise not valid or inactive the key should provide
	//rotation.InvalidTime() as the result of this invocation.
	Created() time.Time
}

type KeyList []Key

//KeyStore abstracts operations to be performed against a key store for rotational capabilities.
type KeyStore interface {
	CreateKey(ctx context.Context) (Key, error)
	DeleteKey(ctx context.Context, key Key) error
	ListKeys(ctx context.Context) (KeyList, error)
	//MaximumKeys is the total number of keys which can be created per user within this store.  This is useful in cases
	//where multiple keys maybe in grace, among others.
	MaximumKeys() int
}
