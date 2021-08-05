//Package rotation provides a set of vendor-agnostic algorithms and utilities for managing key life cycles including
//creation, maintenance, and retirement.
package rotation

import (
	"context"
)

//KeyRotationPlan is the instructions to realize a specific rotation strategy against a KeyStore.
type KeyRotationPlan struct {
	CreateKey   bool
	DestroyKeys KeyList
	goodKeys    KeyList
}

//Apply performs the desired operations against a given store.  If successful a KeyList of healthy keys are returned.
func (plan *KeyRotationPlan) Apply(ctx context.Context, store KeyStore) (KeyList, error) {
	knownKeys := plan.goodKeys
	for _, k := range plan.DestroyKeys {
		if err := store.DeleteKey(ctx, k); err != nil {
			return nil, err
		}
	}
	if plan.CreateKey {
		key, err := store.CreateKey(ctx)
		if err != nil {
			return nil, err
		}
		knownKeys = append(knownKeys, key)
	}
	return knownKeys, nil
}
