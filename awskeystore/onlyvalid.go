package awskeystore

import (
	"context"
	"github.com/truewhitespace/key-rotation/rotation"
)

//OnlyValidKeys is a decorator designed to file AWS access key.  Any keys not within the known set will be invalidated
//as though the keys are expired.
type OnlyValidKeys struct {
	rotation.KeyStoreDecorator
	validKeys StringSlice
}

func (o *OnlyValidKeys) ListKeys(ctx context.Context) (rotation.KeyList, error) {
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
