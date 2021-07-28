package rotation

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func NewKeyRotation(maximumAge time.Duration, graceAge time.Duration) (*KeyRotation, error) {
	if maximumAge <= graceAge {
		return nil, fmt.Errorf("maximum age (%d) must be greater than or equal to grace age (%d)", maximumAge, graceAge)
	}
	return &KeyRotation{
		maximumAge: maximumAge,
		graceAge:   graceAge,
	}, nil
}

type KeyRotation struct {
	maximumAge time.Duration
	graceAge   time.Duration
}

func (k *KeyRotation) Plan(ctx context.Context, store KeyStore) (*KeyRotationPlan, error) {
	graceStart := time.Now().Add(-1 * k.graceAge)
	destroyBefore := time.Now().Add(-1 * k.maximumAge)

	validKeys := make(KeyList, 0)
	graceKeys := make(KeyList, 0)
	expiredKeys := make(KeyList, 0)

	keys, err := store.ListKeys(ctx)
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		created := k.Created()
		if created.Before(destroyBefore) {
			expiredKeys = append(expiredKeys, k)
		} else if created.Before(graceStart) {
			graceKeys = append(graceKeys, k)
		} else {
			validKeys = append(validKeys, k)
		}
	}

	graceKeyCount := len(graceKeys)
	totalKeys := graceKeyCount + len(validKeys)
	willCreate := len(validKeys) == 0
	if willCreate {
		totalKeys++
	}

	if totalKeys >= store.MaximumKeys() {
		if graceKeyCount > 0 {
			//destroy grace key at random
			expiredKeys = append(expiredKeys, graceKeys[0])
		} else {
			return nil, errors.New("no grace keys or available slots")
		}
	}

	return &KeyRotationPlan{
		CreateKey:   willCreate,
		DestroyKeys: expiredKeys,
		goodKeys:    validKeys,
	}, nil
}

type KeyRotationPlan struct {
	CreateKey   bool
	DestroyKeys KeyList
	goodKeys    KeyList
}

//Apply performs the desired operations against the a given store.  If successful a KeyList of healthy keys are returned.
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
