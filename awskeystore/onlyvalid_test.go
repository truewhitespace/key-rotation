package awskeystore

import (
	"context"
	"github.com/truewhitespace/key-rotation/rotation"
	"testing"
	"time"
)

type fakeSource struct {
	AWSUserKeyStore
	returnedKeys rotation.KeyList
}

func (f *fakeSource) ListKeys(ctx context.Context) (rotation.KeyList, error) {
	return f.returnedKeys, nil
}

func newFakeSourceWithKeys(keys rotation.KeyList) rotation.KeyStore {
	return &fakeSource{
		AWSUserKeyStore: AWSUserKeyStore{},
		returnedKeys:    keys,
	}
}

func newTestKeyFilter(validKeys []string, keys rotation.KeyList) rotation.KeyStore {
	source := newFakeSourceWithKeys(keys)
	return NewOnlyValidKeys(source, validKeys)
}

func TestUnknownKeyIsInvalidated(t *testing.T) {
	filter := newTestKeyFilter(nil, rotation.KeyList{&AWSAccessKey{
		ID:      "AKAI-1234",
		Secret:  nil,
		created: time.Now(),
	}})

	ctx := context.Background()
	keys, err := filter.ListKeys(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}

	if len(keys) != 1 {
		t.Errorf("expected 1 keys, got %d keys", len(keys))
		return
	}

	if keys[0].Created() != rotation.InvalidTime() {
		t.Errorf("expected key to be set to invalid date, got %+v", keys[0])
	}
}

func TestKnownKeyIsMaintained(t *testing.T) {
	knownKey := &AWSAccessKey{
		ID:      "AKAI-Gratitude",
		Secret:  nil,
		created: time.Now(),
	}

	filter := newTestKeyFilter([]string{knownKey.ID}, rotation.KeyList{knownKey})

	ctx := context.Background()
	keys, err := filter.ListKeys(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}

	if len(keys) != 1 {
		t.Errorf("expected 1 keys, got %d keys", len(keys))
		return
	}

	returnedKey := keys[0].(*AWSAccessKey)
	if returnedKey != knownKey {
		t.Errorf("expected key to be unmodified, got %+v", keys[0])
	}
}

func TestMixIsTreatedPorperly(t *testing.T) {
	knownKey := &AWSAccessKey{
		ID:      "AKAI-Emptyiness",
		Secret:  nil,
		created: time.Now(),
	}
	unknownKey := &AWSAccessKey{
		ID:      "To all those who stood with me",
		Secret:  nil,
		created: time.Now(),
	}

	filter := newTestKeyFilter([]string{knownKey.ID}, rotation.KeyList{knownKey, unknownKey})

	ctx := context.Background()
	keys, err := filter.ListKeys(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}

	if len(keys) != 2 {
		t.Errorf("expected 1 keys, got %d keys", len(keys))
		return
	}

	keysByID := make(map[string]*AWSAccessKey)
	for _, k := range keys {
		keysByID[k.(*AWSAccessKey).ID] = k.(*AWSAccessKey)
	}
	if keysByID[unknownKey.ID].Created() != rotation.InvalidTime() {
		t.Errorf("exepcted invalidated key got %+v", keys[0])
	}

	goodKey := keysByID[knownKey.ID]
	if goodKey != knownKey {
		t.Errorf("expected key to be unmodified, got %+v", keys[0])
	}
}
