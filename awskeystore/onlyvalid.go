package awskeystore

import (
	"context"
	"github.com/truewhitespace/key-rotation/rotation"
)

//NewOnlyValidKeys creates a new KeyStore only allowing the access key IDs specified in validKeys to be considered
//current.  All other keys will appear to have been expired regardless of status.
func NewOnlyValidKeys(backingStore rotation.KeyStore, validKeys StringSlice) rotation.KeyStore {
	return &onlyValidKeys{
		KeyStoreDecorator: rotation.KeyStoreDecorator{Wrapped: backingStore},
		validKeys:         validKeys,
	}
}

//OnlyValidKeys is a decorator designed to file AWS access key.  Any keys not within the known set will be invalidated
//as though the keys are expired.
type onlyValidKeys struct {
	rotation.KeyStoreDecorator
	validKeys StringSlice
}

func (o *onlyValidKeys) ListKeys(ctx context.Context) (rotation.KeyList, error) {
	keys, err := o.Wrapped.ListKeys(ctx)
	if err != nil {
		return nil, err
	}

	out := make(rotation.KeyList, len(keys))
	for i, k := range keys {
		awsKey := k.(*AWSAccessKey)
		if o.validKeys.Contains(awsKey.ID) {
			out[i] = awsKey
		} else {
			out[i] = &AWSAccessKey{
				ID:      awsKey.ID,
				Secret:  nil,
				created: rotation.InvalidTime(),
			}
		}
	}
	return out, nil
}
