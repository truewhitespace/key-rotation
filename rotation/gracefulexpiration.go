package rotation

import (
	"context"
	"errors"
	"fmt"
	"time"
)

//NewGracefulExpiration instantiates a new key rotation object given the maximum age and grace thresholds provided.
func NewGracefulExpiration(maximumAge time.Duration, graceAge time.Duration) (*GracefulExpiration, error) {
	if maximumAge <= graceAge {
		return nil, fmt.Errorf("maximum age (%d) must be greater than or equal to grace age (%d)", maximumAge, graceAge)
	}
	return &GracefulExpiration{
		maximumAge: maximumAge,
		graceAge:   graceAge,
	}, nil
}

//GracefulExpiration is an algorithm for planning key rotation given a valid key period, and a grace period.  GracefulExpiration will
//attempt to key one key in the active state at all times and destroy any keys exceeding the maximum duration.
//
//If a KeyStore has reached a limit with all keys being in the grace period then one grace key will be selected at
//random to be destroyed.
type GracefulExpiration struct {
	maximumAge time.Duration
	graceAge   time.Duration
}

func (k *GracefulExpiration) Plan(ctx context.Context, store KeyStore) (*KeyRotationPlan, error) {
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
