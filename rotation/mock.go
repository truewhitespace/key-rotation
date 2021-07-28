package rotation

import (
	"context"
	"time"
)

func newMock() *mockKeyStore {
	return &mockKeyStore{
		keys:         make(KeyList, 0),
		createdKey:   false,
		deletedKeys:  nil,
		maximumCount: 1024,
	}
}

func newEmptyMock() KeyStore {
	return newMock()
}

func newMockWithKey() KeyStore {
	store := newMock()
	store.mockGoodKey()
	return store
}

func newMockInGrace() KeyStore {
	store := newMock()
	store.mockInGrace()
	return store
}

func newMockWithExpiredKey() KeyStore {
	store := newMock()
	store.mockExpired(5)
	return store
}

func newMockWithMultipleExpiredKey() KeyStore {
	store := newMock()
	store.mockExpired(15)
	store.mockExpired(5)
	return store
}

type mockKeyStore struct {
	keys         KeyList
	createdKey   bool
	deletedKeys  KeyList
	maximumCount int
}

func (m *mockKeyStore) CreateKey(ctx context.Context) (Key, error) {
	m.createdKey = true
	return &mockKey{created: time.Now()}, nil
}

func (m *mockKeyStore) DeleteKey(ctx context.Context, key Key) error {
	m.deletedKeys = append(m.deletedKeys, key)
	return nil
}

func (m *mockKeyStore) ListKeys(ctx context.Context) (KeyList, error) {
	return m.keys, nil
}

func (m *mockKeyStore) MaximumKeys() int {
	return m.maximumCount
}

func (m *mockKeyStore) appendKeyExpiring(at int) *mockKey {
	if len(m.keys) > m.maximumCount {
		panic("exceeds maximum count")
	}
	key := &mockKey{created: time.Now().Add(-1 * time.Duration(at) * time.Second)}
	m.keys = append(m.keys, key)
	return key
}

func (m *mockKeyStore) mockExpired(secondsAfter int) *mockKey {
	expiry := 60
	return m.appendKeyExpiring(secondsAfter + expiry)
}

func (m *mockKeyStore) mockInGrace() *mockKey {
	return m.appendKeyExpiring(45)
}

func (m *mockKeyStore) mockGoodKey() *mockKey {
	return m.appendKeyExpiring(0)
}

type mockKey struct {
	created time.Time
}

func (m *mockKey) Created() time.Time {
	return m.created
}
